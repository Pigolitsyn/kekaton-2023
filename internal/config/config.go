package config

import "github.com/caarlos0/env/v9"

type Config struct {
	Name   string `env:"APP_NAME"     envDefault:"AwesomeServiceTrustMeBro"`
	Port   string `env:"ADDRESS"      envDefault:":1337"`
	DSN    string `env:"POSTGRES_DSN" envDefault:"postgres://user:password@localhost:5432/db"`
	Secret string `env:"SECRET"       envDefault:"yolo"`
}

func Load() (*Config, error) {
	cfg := &Config{}

	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
