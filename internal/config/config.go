package config

import (
	"github.com/kovetskiy/ko"
	"gopkg.in/yaml.v2"
)

type Database struct {
	DatabaseURI  string `yaml:"database_uri" required:"true" env:"DATABASE_URI"`
	DatabaseName string `yaml:"database_name" required:"true" env:"DATABASE_NAME"`
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
