package conf

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/eviltomorrow/rogue/lib/zlog"
)

type Config struct {
	Path        string `json:"path" toml:"-"`
	ServiceName string `json:"service-name" toml:"service-name"`
	Etcd        Etcd   `json:"etcd" toml:"etcd"`
	Log         Log    `json:"log" toml:"log"`
	SMTPFile    string `json:"smtp-file" toml:"smtp-file"`
}

func (cg *Config) String() string {
	buf, _ := json.Marshal(cg)
	return string(buf)
}

type Etcd struct {
	Endpoints []string `json:"endpoints" toml:"endpoints"`
}

type Log struct {
	DisableTimestamp bool   `json:"disable-timestamp" toml:"disable-timestamp"`
	Level            string `json:"level" toml:"level"`
	Format           string `json:"format" toml:"format"`
	MaxSize          int    `json:"maxsize" toml:"maxsize"`
	FilePath         string `json:"filename" toml:"filename"`
}

func (c *Config) FindAndLoad(path string, override []func(cfg *Config) error) error {
	findPath := func(path string) (string, error) {
		var possibleConf = []string{
			path,
			"../etc/global.conf",
		}
		for _, path := range possibleConf {
			if path == "" {
				continue
			}
			if _, err := os.Stat(path); err == nil {
				fp, err := filepath.Abs(path)
				if err == nil {
					return fp, nil
				}
				return path, nil
			}
		}
		if path == "" {
			return "", nil
		}
		return "", fmt.Errorf("not found conf file, possible conf %v", possibleConf)
	}

	fp, err := findPath(path)
	if err != nil {
		return err
	}

	if fp != "" {
		if _, err := toml.DecodeFile(fp, c); err != nil {
			return err
		}
	}

	for _, f := range override {
		if err := f(c); err != nil {
			return err
		}
	}
	c.Path = fp

	return nil
}

var Global = Config{
	ServiceName: "rogue-email",
	Etcd: Etcd{
		Endpoints: []string{"127.0.0.1:2379"},
	},
	Log: Log{
		DisableTimestamp: false,
		Level:            "info",
		Format:           "text",
		MaxSize:          20,
		FilePath:         "../log/data.log",
	},
	SMTPFile: "smtp.json",
}

func SetupGlobalLog(l Log) error {
	global, prop, err := zlog.InitLogger(&zlog.Config{
		Level:            l.Level,
		Format:           l.Format,
		DisableTimestamp: l.DisableTimestamp,
		File: zlog.FileLogConfig{
			Filename:   l.FilePath,
			MaxSize:    l.MaxSize,
			MaxDays:    30,
			MaxBackups: 30,
			Compress:   true,
		},
		DisableStacktrace:   true,
		DisableErrorVerbose: true,
	})
	if err != nil {
		return err
	}
	zlog.ReplaceGlobals(global, prop)
	return nil
}
