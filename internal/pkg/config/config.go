package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

// PairsConfig represents ...
type PairsConfig struct {
	// Default file
	SourceType  string `env:"PAIRS_SOURCE_TYPE"     envDefault:"file"`
	Source      string `env:"PAIRS_SOURCE,required"`
	PollingRate int    `env:"PAIRS_POLLING_RATE"    envDefault:"30"`
}

type AssetsConfig struct {
	SourceType  string `env:"ASSETS_SOURCE_TYPE"     envDefault:"file"`
	Source      string `env:"ASSETS_SOURCE,required"`
	PollingRate int    `env:"ASSETS_POLLING_RATE"    envDefault:"30"`
}

type ServerConfig struct {
	Addr              string `env:"SERVER_ADDR"                envDefault:":8080"`
	ReadTimeout       int    `env:"SERVER_READ_TIMEOUT"        envDefault:"10"`
	ReadHeaderTimeout int    `env:"SERVER_READ_HEADER_TIMEOUT" envDefault:"10"`
	WriteTimeout      int    `env:"SERVER_WRITE_TIMEOUT"       envDefault:"10"`
}

type Config struct {
	KafkaBrokers    string `env:"KAFKA_BROKERS,required"`
	DBConnString    string `env:"DB_CONN_STRING",required`
	Pairs           PairsConfig
	Assets          AssetsConfig
	Server          ServerConfig
	ShutdownTimeout int `env:"SHUTDOWN_TIMEOUT" envDefault:"10"`
}

func Load() (*Config, error) {
	var config Config

	err := env.Parse(&config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse env: %w", err)
	}

	return &config, nil
}
