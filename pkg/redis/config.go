package redis

import "github.com/spf13/viper"

// config holds the Redis connection settings read from environment variables.
type config struct {
	Address  string `mapstructure:"REDIS_ADDR"`
	Password string `mapstructure:"REDIS_PASSWORD"`
	DB       int    `mapstructure:"REDIS_DB"`
}

// newConfig reads the settings. With a non-empty envPrefix, PREFIX_REDIS_ADDR is
// looked up first and REDIS_ADDR is the fallback (same for the other fields).
func newConfig(envPrefix string) (*config, error) {
	v := viper.New()
	v.SetDefault("REDIS_ADDR", "localhost:6379")
	v.SetDefault("REDIS_PASSWORD", "")
	v.SetDefault("REDIS_DB", 0)

	for _, key := range []string{"REDIS_ADDR", "REDIS_PASSWORD", "REDIS_DB"} {
		envNames := []string{key}
		if envPrefix != "" {
			envNames = []string{envPrefix + "_" + key, key}
		}
		// BindEnv uses the first env name that is set.
		if err := v.BindEnv(append([]string{key}, envNames...)...); err != nil {
			return nil, err
		}
	}

	cfg := &config{}
	err := v.Unmarshal(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
