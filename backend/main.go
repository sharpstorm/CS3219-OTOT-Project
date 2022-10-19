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
	appConfig := util.LoadEnvVariables()

	dbConn, err := database.ConnectDatabase(
		appConfig.DbUrl,
		appConfig.DbUsername,
		appConfig.DbPassword,
		appConfig.DbName,
	)
	if err != nil {
		log.Fatalln("Failed to connect to database")
	}

	tokenAuthenticator := auth.NewTokenAuthenticator(dbConn)

	log.Println("Starting server")
	server := server.CreateHTTPServer(uint16(appConfig.Port))
	attachCardController(server, dbConn, tokenAuthenticator)
	server.AddAssetRoute("/static/*filepath", "./static/static")
	server.AddStaticRoute("/", "./static/index.html")
	server.AddStaticRoute("/favicon.ico", "./static/favicon.ico")
	server.AddStaticRoute("/robots.txt", "./static/robots.txt")

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
