package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type JSONServerConfig struct {
	Name string
	Port string
}

type AppConfig struct {
	JServer1 JSONServerConfig
	JServer2 JSONServerConfig
}

func Load() (*AppConfig, error) {
	viper.AutomaticEnv()

	cfg := &AppConfig{
		JServer1: JSONServerConfig{
			Name: viper.GetString("JSERVER1_NAME"),
			Port: viper.GetString("JSERVER1_PORT"),
		},
		JServer2: JSONServerConfig{
			Name: viper.GetString("JSERVER2_NAME"),
			Port: viper.GetString("JSERVER2_PORT"),
		},
	}

	if cfg.JServer1.Name == "" || cfg.JServer1.Port == "" {
		return nil, fmt.Errorf("missing JSERVER1_* envs")
	}

	if cfg.JServer2.Name == "" || cfg.JServer2.Port == "" {
		return nil, fmt.Errorf("missing JSERVER2_* envs")
	}

	return cfg, nil
}

func (j JSONServerConfig) BaseURL() string {
	return fmt.Sprintf("http://%s:%s", j.Name, j.Port)
}
