package config

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Config is a struct that holds all environment variables
type Config struct {
	// App
	AppIP   string `env:"APP_IP" flag:"app-ip" usage:"ip for app to listen"`
	AppPort string `env:"APP_PORT" flag:"app-port" usage:"port for app to listen"`
	// DB
	DBName     string `env:"DB_NAME" flag:"db-name" usage:"name of database"`
	DBHost     string `env:"DB_HOST" flag:"db-host" usage:"host of database"`
	DBPort     string `env:"DB_PORT" flag:"db-port" usage:"port of database"`
	DBUser     string `env:"DB_USER" flag:"db-user" usage:"user of database"`
	DBPassword string `env:"DB_PASSWORD" flag:"db-password" usage:"password for user in database"`
	// Log
	Mode string `env:"MODE" flag:"mode" usage:"development|production"`
}

// Load parses environment variables and flags, flags have higher priority
// Returns Config and an error if any of the required variables are not set
func Load() (*Config, error) {
	cfg := &Config{}

	v := viper.New()

	cfgType := reflect.TypeOf(*cfg)

	for i := 0; i < cfgType.NumField(); i++ {
		field := cfgType.Field(i)

		envTag := field.Tag.Get("env")
		flagTag := field.Tag.Get("flag")
		usageTag := field.Tag.Get("usage")

		if envTag == "" && flagTag == "" {
			continue
		}

		if envTag != "" {
			_ = v.BindEnv(field.Name, envTag)
		}

		if flagTag != "" {
			pflag.String(flagTag, "", usageTag)
			_ = v.BindPFlag(field.Name, pflag.Lookup(flagTag))
		}

	}

	pflag.Parse()

	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	for i := 0; i < cfgType.NumField(); i++ {
		field := cfgType.Field(i)

		value := v.GetString(field.Name)
		if value == "" {
			return nil, fmt.Errorf("configuration for %s is required but not set", field.Name)
		}
	}

	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("unable to decode into config struct: %w", err)
	}

	return cfg, nil
}
