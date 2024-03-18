package http

import "os"

type Config struct {
	Host string
	Port string
}

func NewConfig() Config {
	return Config{
		Host: os.Getenv("SERVER_HOST"),
		Port: os.Getenv("HTTP_PORT"),
	}
}
