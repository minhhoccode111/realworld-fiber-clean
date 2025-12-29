package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type (
	// Config -.
	Config struct {
		App     App
		HTTP    HTTP
		Log     Log
		PG      PG
		Metrics Metrics
		Swagger Swagger
		JWT     JWT
		CORS    CORS
	}

	// App -.
	App struct {
		Name    string `env:"APP_NAME,required"`
		Version string `env:"APP_VERSION,required"`
	}

	// HTTP -.
	HTTP struct {
		Port           string `env:"HTTP_PORT,required"`
		UsePreforkMode bool   `env:"HTTP_USE_PREFORK_MODE" envDefault:"false"`
	}

	// Log -.
	Log struct {
		Level string `env:"LOG_LEVEL,required"`
	}

	// PG -.
	PG struct {
		PoolMax int    `env:"PG_POOL_MAX,required"`
		URL     string `env:"PG_URL,required"`
	}

	// Metrics -.
	Metrics struct {
		Enabled bool `env:"METRICS_ENABLED" envDefault:"true"`
	}

	// Swagger -.
	Swagger struct {
		Enabled bool `env:"SWAGGER_ENABLED" envDefault:"false"`
	}

	// JWT -.
	JWT struct {
		Issuer     string        `env:"JWT_ISSUER,required"`
		Secret     string        `env:"JWT_SECRET,required"`
		Expiration time.Duration `env:"JWT_EXPIRATION,required"`
	}

	// CORS -.
	CORS struct {
		AllowOrigins     string `env:"CORS_ALLOW_ORIGINS,required"`
		AllowCredentials bool   `env:"CORS_ALLOW_CREDENTIALS,required"`
		AllowHeaders     string `env:"CORS_ALLOW_HEADERS,required"`
		AllowMethods     string `env:"CORS_ALLOW_METHODS,required"`
	}
)

// NewConfig returns app config.
func NewConfig() (*Config, error) {
	// load .env file to make 'air' work for hot reload, without this the app
	// still run normally with 'make run' etc. but not 'air'
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	return cfg, nil
}
