package config

import (
	"gopkg.in/yaml.v2"
	"os"
)

type Config struct {
	Http Http `yaml:"http"`
}

type Http struct {
	Port int `yaml:"port"`
}

func NewConfig(filePath string) (*Config, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
