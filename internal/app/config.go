package app

import (
	"github.com/sivistrukov/vk-assigment/internal/entrypoints/http"
	"github.com/sivistrukov/vk-assigment/internal/infrastructure/postgresql"
)

type Config struct {
	Http     http.Config
	Database postgresql.Config
}

func NewConfig() Config {
	return Config{
		Http:     http.NewConfig(),
		Database: postgresql.NewConfig(),
	}
}
