package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/knadh/koanf/parsers/dotenv"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

// Config holds all application configuration.
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Log      LogConfig
}

// ServerConfig holds HTTP server settings.
type ServerConfig struct {
	Host         string        `koanf:"server_host"`
	Port         int           `koanf:"server_port"`
	ReadTimeout  time.Duration `koanf:"server_read_timeout"`
	WriteTimeout time.Duration `koanf:"server_write_timeout"`
	IdleTimeout  time.Duration `koanf:"server_idle_timeout"`
}

// DatabaseConfig holds database connection settings.
type DatabaseConfig struct {
	Host            string        `koanf:"db_host"`
	Port            int           `koanf:"db_port"`
	User            string        `koanf:"db_user"`
	Password        string        `koanf:"db_password"`
	DBName          string        `koanf:"db_name"`
	SSLMode         string        `koanf:"db_ssl_mode"`
	MaxOpenConns    int           `koanf:"db_max_open_conns"`
	MaxIdleConns    int           `koanf:"db_max_idle_conns"`
	ConnMaxLifetime time.Duration `koanf:"db_conn_max_lifetime"`
}

// DSN returns the PostgreSQL connection string.
func (d DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		d.User, d.Password, d.Host, d.Port, d.DBName, d.SSLMode,
	)
}

// LogConfig holds logging settings.
type LogConfig struct {
	Level  string `koanf:"log_level"`
	Format string `koanf:"log_format"`
}

// Load reads configuration from .env file and environment variables.
// Environment variables take precedence over .env file values.
func Load() (*Config, error) {
	k := koanf.New(".")

	// Load .env file (optional — won't fail if missing)
	if err := k.Load(file.Provider(".env"), dotenv.Parser()); err != nil {
		// Ignore file-not-found; fail on parse errors
		_ = err
	}

	// Load environment variables (override .env)
	if err := k.Load(env.Provider("", ".", func(s string) string {
		return strings.ToLower(s)
	}), nil); err != nil {
		return nil, fmt.Errorf("loading env vars: %w", err)
	}

	cfg := &Config{
		Server: ServerConfig{
			Host:         k.String("server_host"),
			Port:         k.Int("server_port"),
			ReadTimeout:  k.Duration("server_read_timeout"),
			WriteTimeout: k.Duration("server_write_timeout"),
			IdleTimeout:  k.Duration("server_idle_timeout"),
		},
		Database: DatabaseConfig{
			Host:            k.String("db_host"),
			Port:            k.Int("db_port"),
			User:            k.String("db_user"),
			Password:        k.String("db_password"),
			DBName:          k.String("db_name"),
			SSLMode:         k.String("db_ssl_mode"),
			MaxOpenConns:    k.Int("db_max_open_conns"),
			MaxIdleConns:    k.Int("db_max_idle_conns"),
			ConnMaxLifetime: k.Duration("db_conn_max_lifetime"),
		},
		Log: LogConfig{
			Level:  k.String("log_level"),
			Format: k.String("log_format"),
		},
	}

	setDefaults(cfg)

	return cfg, nil
}

func setDefaults(cfg *Config) {
	if cfg.Server.Host == "" {
		cfg.Server.Host = "0.0.0.0"
	}
	if cfg.Server.Port == 0 {
		cfg.Server.Port = 8080
	}
	if cfg.Server.ReadTimeout == 0 {
		cfg.Server.ReadTimeout = 10 * time.Second
	}
	if cfg.Server.WriteTimeout == 0 {
		cfg.Server.WriteTimeout = 10 * time.Second
	}
	if cfg.Server.IdleTimeout == 0 {
		cfg.Server.IdleTimeout = 120 * time.Second
	}
	if cfg.Database.Host == "" {
		cfg.Database.Host = "localhost"
	}
	if cfg.Database.Port == 0 {
		cfg.Database.Port = 5432
	}
	if cfg.Database.User == "" {
		cfg.Database.User = "postgres"
	}
	if cfg.Database.Password == "" {
		cfg.Database.Password = "postgres"
	}
	if cfg.Database.DBName == "" {
		cfg.Database.DBName = "myapp"
	}
	if cfg.Database.SSLMode == "" {
		cfg.Database.SSLMode = "disable"
	}
	if cfg.Database.MaxOpenConns == 0 {
		cfg.Database.MaxOpenConns = 25
	}
	if cfg.Database.MaxIdleConns == 0 {
		cfg.Database.MaxIdleConns = 5
	}
	if cfg.Database.ConnMaxLifetime == 0 {
		cfg.Database.ConnMaxLifetime = 5 * time.Minute
	}
	if cfg.Log.Level == "" {
		cfg.Log.Level = "debug"
	}
	if cfg.Log.Format == "" {
		cfg.Log.Format = "text"
	}
}
