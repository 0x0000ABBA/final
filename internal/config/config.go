package config

import (
	"flag"
	"fmt"
	"os"
)

// Config is a struct that holds all configuration variables
type Config struct {
	// App
	AppIP   string
	AppPort string
	// DB
	DBName     string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	// Log
	Mode string
	// Telemetry
	TelemetryEndpoint string
	// Metrics
	MetricsEndpoint string
}

// Load parses environment variables and flags, flags have higher priority
// Returns Config and an error if any of the required variables are not set
func Load() (*Config, error) {

	appIPFlag := flag.String("app-ip", "", "IP address for the application")
	appPortFlag := flag.String("app-port", "", "Port for the application")

	dbNameFlag := flag.String("db-name", "", "Database name")
	dbHostFlag := flag.String("db-host", "", "Database host")
	dbPortFlag := flag.String("db-port", "", "Database port")
	dbUserFlag := flag.String("db-user", "", "Database user")
	dbPasswordFlag := flag.String("db-password", "", "Database password")

	modeFlag := flag.String("mode", "", "Application mode (e.g., devцццelopment, production)")

	telemetryEndpointFlag := flag.String("telemetry-endpoint", "", "Telemetry endpoint URL")
	metricsEndpointFlag := flag.String("metrics-endpoint", "", "Metrics endpoint URL")

	flag.Parse()

	getValue := func(flagValue *string, envVar string) string {
		if *flagValue != "" {
			return *flagValue
		}
		return os.Getenv(envVar)
	}

	config := &Config{
		AppIP:             getValue(appIPFlag, "APP_IP"),
		AppPort:           getValue(appPortFlag, "APP_PORT"),
		DBName:            getValue(dbNameFlag, "DB_NAME"),
		DBHost:            getValue(dbHostFlag, "DB_HOST"),
		DBPort:            getValue(dbPortFlag, "DB_PORT"),
		DBUser:            getValue(dbUserFlag, "DB_USER"),
		DBPassword:        getValue(dbPasswordFlag, "DB_PASSWORD"),
		Mode:              getValue(modeFlag, "MODE"),
		TelemetryEndpoint: getValue(telemetryEndpointFlag, "TELEMETRY_ENDPOINT"),
		MetricsEndpoint:   getValue(metricsEndpointFlag, "METRICS_ENDPOINT"),
	}

	missingFields := []string{}

	if config.AppIP == "" {
		missingFields = append(missingFields, "AppIP (flag: -app-ip or env: APP_IP)")
	}
	if config.AppPort == "" {
		missingFields = append(missingFields, "AppPort (flag: -app-port or env: APP_PORT)")
	}
	if config.DBName == "" {
		missingFields = append(missingFields, "DBName (flag: -db-name or env: DB_NAME)")
	}
	if config.DBHost == "" {
		missingFields = append(missingFields, "DBHost (flag: -db-host or env: DB_HOST)")
	}
	if config.DBPort == "" {
		missingFields = append(missingFields, "DBPort (flag: -db-port or env: DB_PORT)")
	}
	if config.DBUser == "" {
		missingFields = append(missingFields, "DBUser (flag: -db-user or env: DB_USER)")
	}
	if config.DBPassword == "" {
		missingFields = append(missingFields, "DBPassword (flag: -db-password or env: DB_PASSWORD)")
	}
	if config.Mode == "" {
		missingFields = append(missingFields, "Mode (flag: -mode or env: MODE)")
	}
	if config.TelemetryEndpoint == "" {
		missingFields = append(missingFields, "TelemetryEndpoint (flag: -telemetry-endpoint or env: TELEMETRY_ENDPOINT)")
	}
	if config.MetricsEndpoint == "" {
		missingFields = append(missingFields, "MetricsEndpoint (flag: -metrics-endpoint or env: METRICS_ENDPOINT)")
	}

	if len(missingFields) > 0 {
		return nil, fmt.Errorf("missing required configuration fields: %v", missingFields)
	}

	return config, nil
}
