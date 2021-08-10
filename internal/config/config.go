package config

import "os"

type Config struct {
	Port string

	PostgresAuthenticationData
}

type PostgresAuthenticationData struct {
	PostgresUsername string
	PostgresPassword string
	PostgresDatabase string
	PostgresHost     string
	PostgresPort     string
}

func New() (*Config, error) {
	return &Config{
		Port: os.Getenv("PORT"),
		PostgresAuthenticationData: PostgresAuthenticationData{
			PostgresUsername: os.Getenv("DB_USER"),
			PostgresPassword: os.Getenv("DB_PASSWORD"),
			PostgresHost:     os.Getenv("DB_HOST"),
			PostgresPort:     os.Getenv("DB_PORT"),
			PostgresDatabase: os.Getenv("DB_BASE"),
		},
	}, nil
}
