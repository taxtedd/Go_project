package config

import (
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
	"os"
)

type Config struct {
	Http Http `yaml:"http"`
}

type Http struct {
	Port int `yaml:"port"`
}

func NewConfig(filePath string, logger *zap.Logger) (*Config, error) {
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
