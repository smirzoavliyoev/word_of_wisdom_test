package config

import (
	"os"
	"path"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/sirupsen/logrus"
)

const (
	DefaultConfigPath = ".wisdom-server"
	DefaultConfigFile = "config.yaml"
)

type Config struct {
	Host      string `yaml:"host"`
	Port      uint32 `yaml:"port"`
	Timeout   int64  `yaml:"timeout"`
	SecretKey string `yaml:"secret_key"`

	ConfigPath string
}

type Options func(cfg *Config)

//as second argument provide path
func WithSpecificConfigPathOption(cfg *Config) {
	cfg.ConfigPath = "configs/config.yaml"
}

func ReadConfig(options ...Options) (*Config, error) {

	config := &Config{}

	for _, option := range options {
		option(config)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configPath := config.ConfigPath

	if configPath == "" {
		configPath = path.Join(homeDir, DefaultConfigPath, DefaultConfigFile)
	}

	if err = cleanenv.ReadConfig(configPath, config); err != nil {
		return nil, err
	}

	logrus.WithField("path", configPath).Info("Config loaded")

	return config, nil
}
