package config

import (
	"os"	
)

type Config struct {
	DBDSN string
	Port string
}

func Load() *Config {
	return &Config{
		DBDSN: os.Getenv("DB_DSN"),
		Port:  os.Getenv("PORT"),
	}
}