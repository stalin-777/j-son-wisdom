package app

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

// Config represents configuration settings for the application
type Cfg struct {
	Network      string `env:"NETWORK" envDefault:"tcp"`
	Address      string `env:"ADDRESS" envDefault:":8080"`
	Algorithm    string `env:"ALGORITHM" envDefault:"sha256"`
	Difficulty   int    `env:"DIFFICULTY" envDefault:"5"`
	IsProduction bool   `env:"IS_PRODUCTION" envDefault:"false"`
}

// NewCfg creates a new configuration object
func NewCfg() (*Cfg, error) {
	cfg := &Cfg{}

	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse environments: %v", err)
	}

	return cfg, nil
}
