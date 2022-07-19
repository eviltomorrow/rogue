package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/eviltomorrow/rogue/app/collector/collect"
	"github.com/eviltomorrow/rogue/app/collector/conf"
	"github.com/eviltomorrow/rogue/app/email/pb"
	"github.com/eviltomorrow/rogue/lib/buildinfo"
	"github.com/eviltomorrow/rogue/lib/etcd"
	"github.com/eviltomorrow/rogue/lib/fs"
	"github.com/eviltomorrow/rogue/lib/grpcclient"
	"github.com/eviltomorrow/rogue/lib/mongodb"
	"github.com/eviltomorrow/rogue/lib/pid"
	"github.com/eviltomorrow/rogue/lib/procutil"
	"github.com/eviltomorrow/rogue/lib/runutil"
	"github.com/eviltomorrow/rogue/lib/self"
	"github.com/eviltomorrow/rogue/lib/util"
	"github.com/eviltomorrow/rogue/lib/zlog"
	"github.com/robfig/cron"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"gopkg.in/natefinch/lumberjack.v2"
)

var rootCommand = &cobra.Command{
	Use:   "rogue-collector",
	Short: "Collector service for collect stock trade data",
	Long:  "Collector service for collect stock trade data",
	Run: func(_ *cobra.Command, _ []string) {
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
	rootCommand.Flags().BoolVarP(&version, "version", "v", false, "version about rogue-collector")
	rootCommand.Flags().BoolVarP(&daemon, "daemon", "d", false, "run rogue-collector in daemon mode")
	rootCommand.Flags().StringVarP(&pidFile, "pid", "p", "../var/run/rogue-collector.pid", "rogue-collector's pid path")
	rootCommand.Flags().StringVarP(&config, "config", "c", "", "rogue-collector's pid path")
	rootCommand.Flags().StringVarP(&mode, "mode", "m", "release", "run rogue-collector mode in [release/debug]")
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
	mongodb.DSN = cfg.MongoDB.DSN
	collect.BaseCode = cfg.Collect.CodeList

	if strings.Count(cfg.Collect.RandomWait, ",") == 1 {
		var attrs = strings.Split(cfg.Collect.RandomWait, ",")
		v1, _ := strconv.Atoi(attrs[0])
		v2, _ := strconv.Atoi(attrs[1])
		if v1 < v2 && v1 > 0 && v2 < 100 {
			collect.RandomWait = [2]int{v1, v2}
		}
	}
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

	closeFunc1, err := pid.CreatePidFile(filepath.Join(runutil.ExecutableDir, "../var/run/rogue-collector.pid"))
	if err != nil {
		return err
	}
	self.RegisterClearFuncs(closeFunc1)

	if err := mongodb.Build(); err != nil {
		return err
	}
	self.RegisterClearFuncs(mongodb.Close)

	client, err := etcd.NewClient()
	if err != nil {
		return err
	}
	self.RegisterClearFuncs(client.Close)

	closeFunc2, err := initCrontab()
	if err != nil {
		return err
	}
	self.RegisterClearFuncs(closeFunc2)

	return nil
}

func initCrontab() (func() error, error) {
	var c = cron.New()
	_, err := c.AddFunc(cfg.Collect.Crontab, func() {
		zlog.Info("Sync data slow begin")
		var (
			count int64
			err   error
			begin = time.Now()
		)
		defer func() {
			if err != nil {
				client, closeFunc, err := grpcclient.NewEmail()
				if err != nil {
					zlog.Error("Create email client failure", zap.Error(err))
					return
				}
				defer closeFunc()

				ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
				defer cancel()

				if _, err := client.Send(ctx, &pb.Mail{
					To: []*pb.Contact{
						{Name: "Shepard", Address: "eviltomorrow@163.com"},
					},
					Subject: fmt.Sprintf("同步数据失败-[%s]", time.Now().Format("2006-01-02")),
					Body:    fmt.Sprintf("错误描述, nest error: %v", err),
				}); err != nil {
					zlog.Error("Send email failure, notify [sync data slow]", zap.Error(err))
					return
				}
			} else {
				_ = 1
			}
		}()

		count, err = collect.SyncDataSlow(cfg.Collect.Source)
		if err != nil {
			return
		}
		zlog.Info("Sync data slow complete", zap.Int64("count", count), zap.Duration("cost", time.Since(begin)))
	})
	if err != nil {
		return nil, err
	}
	c.Start()

	return func() error {
		c.Stop()
		return nil
	}, nil
}
