package config

import "github.com/pyr33x/envy"

type Config struct {
	Server Server
}

type Server struct {
	Port int
}

func New() *Config {
	return &Config{
		Server: Server{
			Port: envy.GetInt("PORT", 80),
		},
	}
}
