package server

import (
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

var k = koanf.New(".")

type Config struct {
	ListenAddress string `koanf:"listen_address"`
	DBPath        string `koanf:"db_path"`
	DBType        string `koanf:"db_type"`
	LogLevel      string `koanf:"log_level"`
	ProxyConfig   struct {
		ServerAddress string `koanf:"address"`
		Token         string `koanf:"token"`
	} `koanf:"proxy"`
}

func LoadConfig(path string) (*Config, error) {
	if path != "" {
		if err := k.Load(file.Provider(path), yaml.Parser()); err != nil {
			return nil, err
		}
	}
	c := Config{
		ListenAddress: ":8081",
		DBPath:        "./fence.db",
		DBType:        "badger",
		LogLevel:      "info",
	}
	if err := k.UnmarshalWithConf("", &c, koanf.UnmarshalConf{Tag: "koanf"}); err != nil {
		return nil, err
	}
	return &c, nil
}
