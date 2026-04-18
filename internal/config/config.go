// Package config provides configuration loading and parsing for the Hades service.
// It supports YAML configuration files, environment variables,
// and explicit overrides for sensitive values like database credentials.
package config

import (
	"fmt"
	"os"
	"time"

	wbf "github.com/wb-go/wbf/config"
)

// Config is the root configuration structure aggregating all service settings.
type Config struct {
	Logger  Logger  `mapstructure:"logger"`   // Logger holds logging-related settings.
	Server  Server  `mapstructure:"server"`   // Server holds HTTP server settings.
	Storage Storage `mapstructure:"database"` // Storage holds database connection and migration settings.
}

// Logger configures the application's logging behaviour.
type Logger struct {
	Debug  bool   `mapstructure:"debug_mode"`    // Debug enables debug-level logging when true.
	LogDir string `mapstructure:"log_directory"` // LogDir specifies the directory where log files are written.
}

// Server defines the HTTP server parameters.
type Server struct {
	Port            string        `mapstructure:"port"`             // Port is the TCP port the server listens on (e.g., ":8080").
	ReadTimeout     time.Duration `mapstructure:"read_timeout"`     // ReadTimeout is the maximum duration for reading the entire request.
	WriteTimeout    time.Duration `mapstructure:"write_timeout"`    // WriteTimeout is the maximum duration before timing out writes of the response.
	MaxHeaderBytes  int           `mapstructure:"max_header_bytes"` // MaxHeaderBytes controls the maximum number of bytes the server will read parsing request headers.
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"` // ShutdownTimeout is the grace period allowed for the server to shut down existing connections.
}

// Storage holds PostgreSQL database connection parameters, connection pool settings,
// and Goose migration configuration.
type Storage struct {
	Dialect            string             `mapstructure:"goose_dialect"`              // Dialect is the database dialect used by Goose (e.g., "postgres").
	MigrationsDir      string             `mapstructure:"goose_migrations_directory"` // MigrationsDir is the filesystem path containing SQL migration files.
	Host               string             `mapstructure:"host"`                       // Host is the database server address.
	Port               string             `mapstructure:"port"`                       // Port is the database server port.
	Username           string             `mapstructure:"username"`                   // Username for database authentication.
	Password           string             `mapstructure:"password"`                   // Password for database authentication.
	DBName             string             `mapstructure:"dbname"`                     // DBName is the name of the database to connect to.
	SSLMode            string             `mapstructure:"sslmode"`                    // SSLMode controls SSL connection behaviour (e.g., "disable", "require").
	MaxOpenConns       int                `mapstructure:"max_open_conns"`             // MaxOpenConns sets the maximum number of open connections to the database.
	MaxIdleConns       int                `mapstructure:"max_idle_conns"`             // MaxIdleConns sets the maximum number of idle connections in the pool.
	ConnMaxLifetime    time.Duration      `mapstructure:"conn_max_lifetime"`          // ConnMaxLifetime is the maximum amount of time a connection may be reused.
	QueryRetryStrategy QueryRetryStrategy `mapstructure:"query_retry_strategy"`       // QueryRetryStrategy configures retry behaviour for failed database queries.
}

// QueryRetryStrategy defines the backoff and retry policy for transient database errors.
type QueryRetryStrategy struct {
	Attempts int           `mapstructure:"attempts"` // Attempts is the maximum number of retry attempts.
	Delay    time.Duration `mapstructure:"delay"`    // Delay is the initial delay between retries.
	Backoff  float64       `mapstructure:"backoff"`  // Backoff is the multiplicative factor applied to delay after each retry.
}

// Load reads and parses configuration from the default files (config.yaml and .env),
// merges environment variables, and returns a populated Config struct.
// It returns an error if the configuration files cannot be read or parsed,
// or if unmarshalling fails.
func Load() (Config, error) {

	cfg := wbf.New()

	if err := cfg.LoadConfigFiles("./config.yaml"); err != nil {
		return Config{}, err
	}

	if err := cfg.LoadEnvFiles(".env"); err != nil && !cfg.GetBool("docker") {
		return Config{}, err
	}

	var conf Config

	if err := cfg.Unmarshal(&conf); err != nil {
		return Config{}, fmt.Errorf("unmarshal config: %w", err)
	}

	loadEnvs(&conf)

	return conf, nil

}

// loadEnvs overrides sensitive database credentials from environment variables.
// It is called after loading the YAML and .env files to ensure the final values
// come from the environment.
func loadEnvs(conf *Config) {

	conf.Storage.Username = os.Getenv("DB_USER")
	conf.Storage.Password = os.Getenv("DB_PASSWORD")

}
