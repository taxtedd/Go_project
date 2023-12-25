package config

import (
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

type Config struct {
	Http Http `yaml:"http"`
}

type Http struct {
	Port int `yaml:"port"`
}

func NewConfig(filePath string) (*Config, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Logger init error. %v", err)
	}
	data, err := os.ReadFile(filePath)
	if err != nil {
		logger.Error("Failed to read config file", zap.String("file_path", filePath), zap.Error(err))
		return nil, err
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		logger.Error("Failed to unmarshal YAML", zap.String("file_path", filePath), zap.Error(err))
		return nil, err
	}

	return &cfg, nil
}
