package config

import (
	"strings"

	"github.com/google/uuid"
	"github.com/spf13/viper"
)

type Config struct {
	Port        string `mapstructure:"PORT"`
	ServiceName string `mapstructure:"SERVICE_NAME"`
	InstanceID  string `mapstructure:"INSTANCE_ID"`
}

func Load() (Config, error) {
	v := viper.New()
	v.AutomaticEnv()
	v.SetDefault("PORT", "8080")
	v.SetDefault("SERVICE_NAME", "golink")
	v.SetDefault("INSTANCE_ID", "")

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return Config{}, err
	}

	cfg.InstanceID = strings.TrimSpace(cfg.InstanceID)
	if cfg.InstanceID == "" {
		cfg.InstanceID = uuid.NewString()
	}
	return cfg, nil
}
