package config

import (
	"gopkg.in/yaml.v2"
	"os"
)

type MongoDBConfig struct {
	URI  string `json:"uri"`
	Name string `json:"name"`
}

type Config struct {
	Port     string        `json:"port"`
	Database MongoDBConfig `json:"database"`
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
