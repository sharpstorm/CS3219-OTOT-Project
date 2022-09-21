package main

import (
	"log"

	"backend.cs3219.comp.nus.edu.sg/auth"
	"backend.cs3219.comp.nus.edu.sg/controller"
	"backend.cs3219.comp.nus.edu.sg/database"
	"backend.cs3219.comp.nus.edu.sg/server"
	"backend.cs3219.comp.nus.edu.sg/util"
)

func main() {
	dbConfig := util.LoadEnvVariables()

	dbConn, err := database.ConnectDatabase(dbConfig.Url, dbConfig.Username, dbConfig.Password, dbConfig.Name)
	if err != nil {
		log.Fatalln("Failed to connect to database")
	}

	tokenAuthenticator := auth.NewTokenAuthenticator(dbConn)

	log.Println("Starting server")
	server := server.CreateHTTPServer(8080)
	attachCardController(server, dbConn, tokenAuthenticator)
	// server.AddStaticRoute("/", "./static/index.html")
	// server.AddAssetRoute("/static/*filepath", "./static/")

	server.Start()
}

func attachCardController(
	server server.HTTPServer,
	dbConnection *database.DatabaseConnection,
	tokenAuthenticator auth.TokenAuthenticator,
) {
	controller := controller.NewCardController(dbConnection, tokenAuthenticator)
	controller.Attach(server)
}
