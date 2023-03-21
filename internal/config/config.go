package config

import "time"

type Config struct {
	RunAddress           string        `env:"RUN_ADDRESS"`
	DatabaseURI          string        `env:"DATABASE_URI"`
	AccrualSystemAddress string        `env:"ACCRUAL_SYSTEM_ADDRESS"`
	AccrualPollInterval  time.Duration `env:"ACCRUAL_POLL_INTERVAL"`
}
