package config

import (
	"errors"
	"flag"
	"fmt"

	"github.com/caarlos0/env/v6"
)

var (
	ErrRunAddressIsNotSet           = errors.New("servers's addres is not set")
	ErrDatabaseURIIsNotSet          = errors.New("servers's database's URI is not set")
	ErrAccrualSystemAddressIsNotSet = errors.New("accrual system's address is not set")
)

type Config struct {
	RunAddress     string `env:"RUN_ADDRESS"`
	DatabaseURI    string `env:"DATABASE_URI"`
	AccrualAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
}

func Make() (*Config, error) {
	config := Config{}

	flag.StringVar(&config.RunAddress, "a", "", "Server's host:port")
	flag.StringVar(&config.DatabaseURI, "d", "", "Database uri")
	flag.StringVar(&config.AccrualAddress, "r", "", "Accrual system's address")
	flag.Parse()

	if err := env.Parse(&config); err != nil {
		return nil, fmt.Errorf("parse env, err=%w", err)
	}

	var err error

	if config.RunAddress == "" {
		err = errors.Join(err, ErrRunAddressIsNotSet)
	}

	if config.DatabaseURI == "" {
		err = errors.Join(err, ErrDatabaseURIIsNotSet)
	}

	if config.AccrualAddress == "" {
		err = errors.Join(err, ErrAccrualSystemAddressIsNotSet)
	}

	if err != nil {
		return nil, fmt.Errorf("bad config, err=%w", err)
	}

	return &config, nil
}
