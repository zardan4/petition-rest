package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	petitions "github.com/zardan4/petition-rest/internal/core"
	"github.com/zardan4/petition-rest/internal/service"
	repository "github.com/zardan4/petition-rest/internal/storage/psql"
	mq_client "github.com/zardan4/petition-rest/internal/transport/mq"
	handlers "github.com/zardan4/petition-rest/internal/transport/rest"
	"github.com/zardan4/petition-rest/pkg/hashing"

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

// @title Petitions REST API Documentation
// @version 1.0
// @description Can be used for writing small petitions interfaces

// @host localhost:8080
// @BasePath  /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	var dbPort string
	if viper.GetString("docker") == "true" {
		dbPort = os.Getenv("DB_INSIDE_PORT")
	} else {
		dbPort = os.Getenv("DB_PORT")
	}

	db, err := repository.NewPostgresDB(repository.PostgresConfig{
		Host: os.Getenv("DB_HOST"),
		Port: dbPort,

		Username: os.Getenv("DB_USERNAME"),
		DBName:   os.Getenv("DB_DBNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	})
	if err != nil {
		logrus.Fatalf("Failed to initialize database: %v", err)
	}

	hasher := hashing.NewSHA256Hasher(os.Getenv("HASHER_SALT"))

	auditClient, err := mq_client.NewClient(os.Getenv("RABBITMQ_USER"), os.Getenv("RABBITMQ_PASSWORD"))
	if err != nil {
		logrus.Fatal("error creating grpc audit client: ", err)
	}

	repo := repository.NewRepository(db)
	services := service.NewService(repo, hasher, []byte(os.Getenv("JWT_SIGNING_KEY")), viper.GetDuration("auth.jwtTTL"), auditClient)
	handlers := handlers.NewHandler(services)

	srv := new(petitions.Server)

	go func() {
		if err := srv.Run(os.Getenv("SERVER_PORT"), handlers.InitRoutes()); err != nil {
			logrus.Fatalf("Failed to initialize routes: %v", err)
		}
	}()

	logrus.Println("server started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Println("server shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logrus.Errorf("error occured on server shutting down: %s", err.Error())
	}

	if err := db.Close(); err != nil {
		logrus.Errorf("error occured on db connection close: %s", err.Error())
	}

	<-ctx.Done()

	logrus.Println("server exiting")
}
