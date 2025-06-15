package main

import (
	"os"

	"gopkg.in/yaml.v3"
)

type DomainConfig struct {
	Domain     string   `yaml:"domain"`
	IP         string   `yaml:"ip"`
	Port       int      `yaml:"port"`
	AllowedIPs []string `yaml:"allowed_ips,omitempty"`
}

type Config struct {
	Listen     string         `yaml:"listen"`
	Domains    []DomainConfig `yaml:"domains"`
	AllowedIPs []string       `yaml:"allowed_ips,omitempty"`
}

func loadConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var cfg Config
	decoder := yaml.NewDecoder(f)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
