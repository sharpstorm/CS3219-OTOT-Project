package util

import (
	"log"
	"os"
	"strconv"
)

type AppConfig struct {
	DbUsername string
	DbPassword string
	DbUrl      string
	DbName     string

	Port int
}

func LoadEnvVariables() AppConfig {
	return loadDatabaseConfig()
}

func loadDatabaseConfig() AppConfig {
	databaseUrl, urlFound := os.LookupEnv("DATABASE_URL")
	databaseUsername, usernameFound := os.LookupEnv("DATABASE_USERNAME")
	databasePassword, passwordFound := os.LookupEnv("DATABASE_PASSWORD")
	databaseName, dbNameFound := os.LookupEnv("DATABASE_NAME")
	portConfig, portFound := os.LookupEnv("APP_PORT")

	if !urlFound || !usernameFound || !passwordFound || !dbNameFound {
		log.Fatal("Database configuration not set.")
	}

	var port int = 80
	if portFound {
		portInt, err := strconv.Atoi(portConfig)
		if err != nil {
			log.Fatal("Port must be an integer")
		}
		port = portInt
	}

	return AppConfig{
		DbUsername: databaseUsername,
		DbPassword: databasePassword,
		DbUrl:      databaseUrl,
		DbName:     databaseName,
		Port:       port,
	}
}
