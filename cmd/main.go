package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	petitions "github.com/zardan4/petition-rest"
	"github.com/zardan4/petition-rest/pkg/handlers"
	"github.com/zardan4/petition-rest/pkg/repository"
	"github.com/zardan4/petition-rest/pkg/service"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// init configs
func init() {
	viper.AddConfigPath("configs")
	// global cfg
	viper.SetConfigName("global")
	if err := viper.ReadInConfig(); err != nil {
		logrus.Fatalf("Failed to initialize configuration: %v", err)
	}

	// .env config
	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("Error loading env file: %v", err)
	}
}

// mode setting
func init() {
	switch viper.GetString("mode") {
	case "release":
		gin.SetMode(gin.ReleaseMode)
	case "debug":
		gin.SetMode(gin.DebugMode)
	case "test":
		gin.SetMode(gin.TestMode)
	default:
		panic(fmt.Sprintf("no such mode {%s} found. set mode in config: release OR debug OR test", viper.GetString("mode")))
	}
}

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	db, err := repository.NewPostgresDB(repository.PostgresConfig{
		Host: os.Getenv("POSTGRES_HOST"),
		Port: os.Getenv("POSTGRES_INSIDE_PORT"),

		Username: os.Getenv("POSTGRES_USERNAME"),
		DBName:   os.Getenv("POSTGRES_DBNAME"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		SSLMode:  os.Getenv("POSTGRES_SSLMODE"),
	})
	if err != nil {
		logrus.Fatalf("Failed to initialize database: %v", err)
	}

	repo := repository.NewRepository(db)
	services := service.NewService(repo)
	handlers := handlers.NewHandler(services)

	srv := new(petitions.Server)

	if err := srv.Run(os.Getenv("SERVER_PORT"), handlers.InitRoutes()); err != nil {
		logrus.Fatalf("Failed to initialize routes: %v", err)
	}

	logrus.Print("server started")
}
