package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Server    ServerConfig `yaml:"http_server"`
	Access    AccessConfig `yaml:"access"`
	DebugMode bool         `yaml:"debug_mode"`
}

type ServerConfig struct {
	Port    string   `yaml:"port"`
	Host    string   `yaml:"host"`
	Origins []string `yaml:"origins"`
}

type AccessConfig struct {
	Token      string   `yaml:"token"`
	AllowedIPs []string `yaml:"allowed_ips"`
}

func Load() (*Config, error) {
	data, err := os.ReadFile("config.yml")
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
