package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/eviltomorrow/rogue/app/email/conf"
	"github.com/eviltomorrow/rogue/app/email/server"
	"github.com/eviltomorrow/rogue/lib/buildinfo"
	"github.com/eviltomorrow/rogue/lib/etcd"
	"github.com/eviltomorrow/rogue/lib/fs"
	"github.com/eviltomorrow/rogue/lib/grpclb"
	"github.com/eviltomorrow/rogue/lib/pid"
	"github.com/eviltomorrow/rogue/lib/procutil"
	"github.com/eviltomorrow/rogue/lib/runutil"
	"github.com/eviltomorrow/rogue/lib/self"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/resolver"
	"gopkg.in/natefinch/lumberjack.v2"
)

var rootCommand = &cobra.Command{
	Use:   "rogue-email",
	Short: "Email service for rogue(common)",
	Long:  "Email service for rogue(common)",
	Run: func(cmd *cobra.Command, args []string) {
		if version {
			fmt.Println(buildinfo.GetVersion())
			os.Exit(0)
		}
		if daemon {
			var writer = &lumberjack.Logger{
				Filename:   filepath.Join(runutil.ExecutableDir, "../log/error.log"),
				MaxSize:    20,
				MaxBackups: 10,
				MaxAge:     28,
				Compress:   true,
			}
			if err := procutil.RunInBackground(runutil.ExecutableName, []string{"-pid", pidFile}, nil, writer); err != nil {
				log.Printf("[F] Run app in background failure, nest error: %v", err)
				os.Exit(1)
			}
			os.Exit(0)
		}

		if err := run(); err != nil {
			log.Printf("[F] Run app failure, nest error: %v", err)
			os.Exit(1)
		}
	},
}

var (
	cfg     = conf.Global
	config  string
	version bool
	daemon  bool
	pidFile string
)

func init() {
	rootCommand.Flags().BoolVarP(&version, "version", "v", false, "version about rogue-email")
	rootCommand.Flags().BoolVarP(&daemon, "daemon", "d", false, "run rogue-email in daemon mode")
	rootCommand.Flags().StringVarP(&pidFile, "pid", "p", "../var/run/rogue-email.pid", "rouge-email's pid path")
	rootCommand.Flags().StringVarP(&config, "config", "c", "", "rouge-email's pid path")
}

func Execute() error {
	rootCommand.CompletionOptions.HiddenDefaultCmd = true

	return rootCommand.Execute()
}

func run() error {
	defer func() {
		for _, err := range self.RunClearFuncs() {
			log.Printf("[E] clear funcs: %v", err)
		}
	}()

	if err := cfg.FindAndLoad(config, nil); err != nil {
		return err
	}
	if err := setupConfig(); err != nil {
		return err
	}
	setupGlobalVars()

	if err := setupRuntime(); err != nil {
		return err
	}

	client, err := etcd.NewClient()
	if err != nil {
		return err
	}
	self.RegisterClearFuncs(client.Close)

	resolver.Register(grpclb.NewBuilder(client))
	var s = &server.GRPC{
		Client: client,
	}
	if s.StartupGRPC(); err != nil {
		return err
	}
	self.RegisterClearFuncs(s.ShutdownGRPC)

	procutil.WaitForSigterm()
	return nil
}

func setupGlobalVars() {
	procutil.HomeDir = runutil.ExecutableDir
	self.ServiceName = "rogue-email"
	etcd.Endpoints = cfg.Etcd.Endpoints
}

func setupConfig() error {
	if err := conf.SetupGlobalLog(cfg.Log); err != nil {
		return err
	}
	return nil
}

func setupRuntime() error {
	for _, dir := range []string{
		filepath.Join(runutil.ExecutableDir, "../log"),
		filepath.Join(runutil.ExecutableDir, "../var/run")} {
		if err := fs.CreateDir(dir); err != nil {
			return err
		}
	}

	closeFunc, err := pid.CreatePidFile(filepath.Join(runutil.ExecutableDir, "../var/run/rouge-email.pid"))
	if err != nil {
		return err
	}
	self.RegisterClearFuncs(closeFunc)

	return nil
}
