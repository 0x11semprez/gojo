package app

import "os"

type Config struct {
	DatabaseUrl string
	Port        string
}

func NewConfig() *Config {
	return &Config{
		DatabaseUrl: os.Getenv("DATABASE_URL"),
		Port:        os.Getenv("PORT"),
	}
}
