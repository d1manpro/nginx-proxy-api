package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Server           ServerConfig `yaml:"http_server"`
	Access           AccessConfig `yaml:"access"`
	Cloudflare       Cloudflare   `yaml:"cloudflare"`
	Email            string       `yaml:"email"`
	NginxCfgTemplate string
	DebugMode        bool `yaml:"debug_mode"`
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

type Cloudflare struct {
	Token   string            `yaml:"token"`
	NodeIP  string            `yaml:"node_ip"`
	Domains map[string]string `yaml:"domains"`
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

	text, err := os.ReadFile("template.conf")
	if err != nil {
		return nil, err
	}
	cfg.NginxCfgTemplate = string(text)

	return &cfg, nil
}
