package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Database DatabaseConfig `mapstructure:"database"`
	Jwt      JwtConfig      `mapstructure:"jwt"`
}

// AppConfig app specific config
type AppConfig struct {
	Env              string `mapstructure:"env"`
	DataPath         string `mapstructure:"data_path"`
	RequestBodyLimit string `mapstructure:"request_body_limit"`
	DateFormat       string
	TimestampFormat  string
	Host             string        `mapstructure:"host"`
	ReadTimeout      time.Duration `mapstructure:"read_timeout"`
	WriteTimeout     time.Duration `mapstructure:"write_timeout"`
	IdleTimeout      time.Duration `mapstructure:"idle_timeout"`
	ContextTimeout   time.Duration `mapstructure:"context_timeout"`
	Port             int           `mapstructure:"port"`
	MaxPageSize      int           `mapstructure:"max_page_size"`
	DefaultPageSize  int           `mapstructure:"default_page_size"`
	RegistrationOpen bool          `mapstructure:"registration_open"`
}

// DatabaseConfig DB specific config
type DatabaseConfig struct {
	Type        string        `mapstructure:"type"`
	Host        string        `mapstructure:"host"`
	Name        string        `mapstructure:"name"`
	Username    string        `mapstructure:"username"`
	Password    string        `mapstructure:"password"`
	SslMode     string        `mapstructure:"ssl_mode"`
	MaxLifeTime time.Duration `mapstructure:"max_life_time"`
	Port        int           `mapstructure:"port"`
	MaxOpenConn int           `mapstructure:"max_open_conn"`
	MaxIdleConn int           `mapstructure:"max_idle_conn"`
	Debug       bool          `mapstructure:"debug"`
}

type JwtConfig struct {
	SecretKey  string        `mapstructure:"secret_key"`
	ExpireTime time.Duration `mapstructure:"expire_time"`
}

// c is the configuration instance
var c Config //nolint:gochecknoglobals

// Get returns all configurations
func Get() Config {
	return c
}

// Load the config
func Load(path string) error {
	viper.SetConfigType("yaml")
	viper.SetConfigFile(path)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config: %w", err)
	}

	if err := viper.Unmarshal(&c); err != nil {
		return fmt.Errorf("failed to unmarshal consul config: %w", err)
	}

	if c.App.RequestBodyLimit == "" {
		c.App.RequestBodyLimit = "20M"
	}

	if c.App.MaxPageSize <= 5 {
		c.App.MaxPageSize = 50
	}

	if c.App.DefaultPageSize <= 5 {
		c.App.DefaultPageSize = 30
	}

	// yyyy-mm-dd
	c.App.DateFormat = "2006-01-02"
	c.App.TimestampFormat = "2006-01-02T15:04:05-0700"

	dataPath := strings.TrimSpace(c.App.DataPath)
	if dataPath == "" {
		dataPath = "."
	}

	c.App.DataPath = dataPath

	return nil
}
