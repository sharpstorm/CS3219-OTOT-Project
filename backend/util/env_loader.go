package util

import (
	"log"
	"os"
)

type databaseConfig struct {
	Username string
	Password string
	Url      string
	Name     string
}

func LoadEnvVariables() databaseConfig {
	return loadDatabaseConfig()
}

func loadDatabaseConfig() databaseConfig {
	databaseUrl, urlFound := os.LookupEnv("DATABASE_URL")
	databaseUsername, usernameFound := os.LookupEnv("DATABASE_USERNAME")
	databasePassword, passwordFound := os.LookupEnv("DATABASE_PASSWORD")
	databaseName, dbNameFound := os.LookupEnv("DATABASE_NAME")

	if !urlFound || !usernameFound || !passwordFound || !dbNameFound {
		log.Fatal("Database configuration not set.")
	}

	return databaseConfig{
		Username: databaseUsername,
		Password: databasePassword,
		Url:      databaseUrl,
		Name:     databaseName,
	}
}
