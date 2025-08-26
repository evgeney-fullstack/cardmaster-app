package main

import (
	"os"

	"github.com/evgeney-fullstack/cardmaster-app/internal/app/handler"
	"github.com/evgeney-fullstack/cardmaster-app/internal/app/repository/mongodb"
	"github.com/evgeney-fullstack/cardmaster-app/internal/app/repository/redis"
	"github.com/evgeney-fullstack/cardmaster-app/internal/app/server"
	"github.com/evgeney-fullstack/cardmaster-app/internal/app/service"
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

	// Initializing a connection to MongoDB using parameters from environment variables
	mdb, err := mongodb.NewMongoDB(mongodb.Config{
		User:     os.Getenv("MONGO_INITDB_ROOT_USERNAME"),
		Password: os.Getenv("MONGO_INITDB_ROOT_PASSWORD"),
		Host:     os.Getenv("MONGO_INITDB_HOST"),
		Port:     os.Getenv("MONGO_INITDB_PORT"),
	})
	if err != nil {
		logrus.Fatalf("failed to initialize db: %s", err.Error())
	}

	// Initializing the connection to Redis
	rdb, err := redis.NewRedisDB(redis.Config{
		Addr:     os.Getenv("REDIS_INITDB_HOST"),
		Port:     os.Getenv("REDIS_INITDB_PORT"),
		Password: os.Getenv("REDIS_INITDB_ROOT_PASSWORD"),
		DB:       os.Getenv("REDIS_INITDB_DB"),
	})
	if err != nil {
		logrus.Fatalf("failed to initialize redis: %s", err.Error())
	}
	defer rdb.Close()

	// Initializing repositories for working with data
	// repos provides access to MongoDB data
	repos := mongodb.NewRepository(mdb)
	// redisRepo provides a caching layer to improve performance
	redisRepo := redis.NewCacheRepository(rdb)

	// Creating a service layer with dependency injection
	// service encapsulates the business logic of the application
	service := service.NewService(repos, redisRepo)

	// Initialization of HTTP handlers with the introduction of a service layer
	// Handlers will use business logic via service
	handlers := handler.NewHandler(service)

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
