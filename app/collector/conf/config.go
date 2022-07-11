package conf

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Config struct {
	ServiceName string      `json:"service-name" toml:"service-name"`
	Etcd        Etcd        `json:"etcd" toml:"etcd"`
	MongoDB     MongoDB     `json:"mongodb" toml:"mongodb"`
	GatherTrade GatherTrade `json:"gather-trade" toml:"gather-trade"`
	Log         Log         `json:"log" toml:"log"`
}

func (cg *Config) String() string {
	buf, _ := json.Marshal(cg)
	return string(buf)
}

type Etcd struct {
	Endpoints []string `json:"endpoints" toml:"endpoints"`
}

type MongoDB struct {
	DSN string `json:"dsn" toml:"dsn"`
}

type GatherTrade struct {
	Source   string   `json:"source" toml:"source"`
	CodeList []string `json:"code-list" toml:"code-list"`
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
	return nil
}

var Global = Config{
	ServiceName: "rogue-email",
	Etcd: Etcd{
		Endpoints: []string{"127.0.0.1:2379"},
	},
	MongoDB: MongoDB{
		DSN: "mongodb://127.0.0.1:27017",
	},
	GatherTrade: GatherTrade{
		Source: "net126",
		CodeList: []string{
			"sh688***",
			"sh605***",
			"sh603***",
			"sh601***",
			"sh600***",
			"sz300***",
			"sz0030**",
			"sz002***",
			"sz001**",
			"sz000***",
		},
	},
	Log: Log{
		DisableTimestamp: false,
		Level:            "info",
		Format:           "text",
		MaxSize:          20,
		FilePath:         "../log/data.log",
	},
}
