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
	"github.com/eviltomorrow/rogue/lib/util"
	"github.com/eviltomorrow/rogue/lib/zlog"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"google.golang.org/grpc/resolver"
	"gopkg.in/natefinch/lumberjack.v2"

	"net/http"
	_ "net/http/pprof"
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
			if err := procutil.RunInBackground(runutil.ExecutableName, []string{"--pid", pidFile}, nil, writer); err != nil {
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
	mode    string
)

func init() {
	rootCommand.Flags().BoolVarP(&version, "version", "v", false, "version about rogue-email")
	rootCommand.Flags().BoolVarP(&daemon, "daemon", "d", false, "run rogue-email in daemon mode")
	rootCommand.Flags().StringVarP(&pidFile, "pid", "p", "../var/run/rogue-email.pid", "rogue-email's pid path")
	rootCommand.Flags().StringVarP(&config, "config", "c", "", "rogue-email's pid path")
	rootCommand.Flags().StringVarP(&mode, "mode", "m", "release", "run rogue-email mode in [release/debug]")
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
	zlog.Info("Config info", zap.String("global", cfg.String()))

	setupGlobalVars()

	if err := setupRuntime(); err != nil {
		return err
	}

	procutil.WaitForSigterm()
	return nil
}

func setupGlobalVars() {
	procutil.HomeDir = runutil.ExecutableDir
	self.ServiceName = cfg.ServiceName
	etcd.Endpoints = cfg.Etcd.Endpoints
}

func setupConfig() error {
	if err := conf.SetupGlobalLog(cfg.Log); err != nil {
		return err
	}
	return nil
}

func setupRuntime() error {
	if mode == "debug" {
		go func() {
			port, err := util.GetAvailablePort()
			if err != nil {
				log.Printf("[F] Http pprof start failure, nest error: %v", err)
				os.Exit(1)
			}

			log.Printf("[I] Http pprof has listened on [http://127.0.0.1:%d/debug/pprof]", port)
			if err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), nil); err != nil {
				log.Fatalf("[F] Listen and serve pprof failure, nest error: %v\r\n", err)
			}
		}()
	}

	for _, dir := range []string{
		filepath.Join(runutil.ExecutableDir, "../log"),
		filepath.Join(runutil.ExecutableDir, "../var/run")} {
		if err := fs.CreateDir(dir); err != nil {
			return err
		}
	}

	closeFunc, err := pid.CreatePidFile(filepath.Join(runutil.ExecutableDir, "../var/run/rogue-email.pid"))
	if err != nil {
		return err
	}
	self.RegisterClearFuncs(closeFunc)

	client, err := etcd.NewClient()
	if err != nil {
		return err
	}
	self.RegisterClearFuncs(client.Close)

	smtp, err := conf.FindSMTP(filepath.Join(filepath.Dir(cfg.Path), cfg.SMTPFile))
	if err != nil {
		return err
	}

	resolver.Register(grpclb.NewBuilder(client))
	var s = &server.GRPC{
		Client: client,
		SMTP:   smtp,
	}
	if s.StartupGRPC(); err != nil {
		return err
	}
	self.RegisterClearFuncs(s.ShutdownGRPC)

	return nil
}
