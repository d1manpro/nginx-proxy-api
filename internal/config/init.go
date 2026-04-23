package config

import (
	"errors"
	"fmt"
	"net"
	"os"
	"regexp"

	"gopkg.in/yaml.v2"
)

var domainRegex = regexp.MustCompile(`^(?i)[a-z0-9-]+(\.[a-z0-9-]+)+$`)

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
	Token       string `yaml:"token"`
	Type        string
	NodeAddress string            `yaml:"node_address"`
	Domains     map[string]string `yaml:"domains"`
}

func Load() (*Config, error) {
	cfgPath := os.Getenv("NPA_CONFIG")
	if cfgPath == "" {
		cfgPath = "config.yml"
	}

	tmplPath := os.Getenv("NPA_TEMPLATE")
	if tmplPath == "" {
		tmplPath = "template.conf"
	}

	data, err := os.ReadFile(cfgPath)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	text, err := os.ReadFile(tmplPath)
	if err != nil {
		return nil, err
	}
	cfg.NginxCfgTemplate = string(text)

	if cfg.Cloudflare.Token == "your_cloudflare_api_token" || cfg.Cloudflare.NodeAddress == "0.0.0.0" {
		return nil, fmt.Errorf("You need to configure Cloudflare API in %s", cfgPath)
	}

	if cfg.Email == "admin@example.com" {
		return nil, errors.New("You need to edit email")
	}

	if net.ParseIP(cfg.Cloudflare.NodeAddress) != nil {
		cfg.Cloudflare.Type = "A"
	} else {
		if ok := domainRegex.MatchString(cfg.Cloudflare.NodeAddress); !ok {
			return nil, fmt.Errorf("node_address is not valid IPv4 or domain: %s", cfg.Cloudflare.NodeAddress)
		}
		cfg.Cloudflare.Type = "CNAME"
	}

	return &cfg, nil
}
