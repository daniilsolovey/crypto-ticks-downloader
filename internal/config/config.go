package config

import (
	"github.com/kovetskiy/ko"
	"gopkg.in/yaml.v2"
)

type Database struct {
	Name      string `yaml:"name" required:"true" env:"DATABASE_NAME"`
	Host      string `yaml:"host" required:"true" env:"DATABASE_HOST"`
	Port      int    `yaml:"port" required:"true" env:"DATABASE_PORT"`
	User      string `yaml:"user" required:"true"`
	Password  string `yaml:"password" required:"true"`
	TableName string `yaml:"table_name" required:"true"`
}

type Config struct {
	WebsocketURL string   `yaml:"websocket_url" required:"true"`
	Database     Database `yaml:"database" required:"true"`
}

func Load(path string) (*Config, error) {
	config := &Config{}
	err := ko.Load(path, config, ko.RequireFile(false), yaml.Unmarshal)
	if err != nil {
		return nil, err
	}

	return config, nil
}
