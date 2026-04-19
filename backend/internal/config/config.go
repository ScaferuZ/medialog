package config

import (
	"log"
	"strconv"

	"github.com/spf13/viper"
)

type Config struct {
	ServerPort  string
	DatabaseURL string
	JWTSecret   string
	TMDBAPIKey  string
}

func Load() Config {
	viper.SetDefault("SERVER_PORT", "8080")
	viper.SetDefault("DATABASE_URL", "")
	viper.SetDefault("JWT_SECRET", "change-me")
	viper.SetDefault("TMDB_API_KEY", "")

	viper.AutomaticEnv()

	cfg := Config{
		ServerPort:  viper.GetString("SERVER_PORT"),
		DatabaseURL: viper.GetString("DATABASE_URL"),
		JWTSecret:   viper.GetString("JWT_SECRET"),
		TMDBAPIKey:  viper.GetString("TMDB_API_KEY"),
	}

	if _, err := strconv.Atoi(cfg.ServerPort); err != nil {
		log.Printf("invalid SERVER_PORT %q, using 8080", cfg.ServerPort)
		cfg.ServerPort = "8080"
	}

	return cfg
}
