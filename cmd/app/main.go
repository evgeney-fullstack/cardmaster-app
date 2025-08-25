package main

import (
	"os"

	"github.com/evgeney-fullstack/cardmaster-app/internal/app/handler"
	"github.com/evgeney-fullstack/cardmaster-app/internal/app/server"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	// Configuring the logs format in JSON for better structuring and compatibility
	// with monitoring systems (Kibana, Elasticsearch, etc.)
	logrus.SetFormatter(new(logrus.JSONFormatter))

	// Loading environment variables from the config.env file
	if err := godotenv.Load("config.env"); err != nil {
		logrus.Fatalf("error loading env variables: %s", err.Error())
	}

	// Initialization of HTTP request handlers
	handlers := handler.NewHandler()

	// Creating a server instance
	srv := new(server.Server)

	// Launching an HTTPS server with configuration from environment variables
	// Using HOST and HOST_PORT from config.env and checking that the variables are not empty
	host := os.Getenv("HOST")
	port := os.Getenv("HOST_PORT")
	if host != "" || port != "" {

		if err := srv.Run(host, port, handlers.InitRoutes()); err != nil {
			logrus.Fatalf("error occurred while running http server: %s", err.Error())
		}

	} else {
		logrus.Fatal("HOST or HOST_PORT environment variables are not set")
	}

}
