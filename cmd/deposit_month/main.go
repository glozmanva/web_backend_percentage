package main

import (
	"os"

	"deposit_month/internal/app/config"
	"deposit_month/internal/app/dsn"
	"deposit_month/internal/app/handler"
	"deposit_month/internal/app/repository"
	"deposit_month/internal/pkg"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	router := gin.Default()

	conf, err := config.NewConfig()
	if err != nil {
		logrus.Fatalf("error loading config: %v", err)
	}

	postgresString := dsn.FromEnv()
	if postgresString == "" {
		logrus.Fatal(
			"database environment variables are not configured",
		)
	}

	rep, err := repository.New(
		&repository.RepositorySettings{
			PostgresDSN:     postgresString,
			MinioEndpoint:   os.Getenv("MINIO_ENDPOINT"),
			MinioAccessKey:  os.Getenv("MINIO_ACCESS_KEY"),
			MinioSecretKey:  os.Getenv("MINIO_SECRET_KEY"),
			MinioBucketName: os.Getenv("MINIO_BUCKET_NAME"),
		},
	)
	if err != nil {
		logrus.Fatalf(
			"error initializing repository: %v",
			err,
		)
	}

	hand := handler.NewHandler(rep)

	application := pkg.NewApp(
		conf,
		router,
		hand,
	)

	application.RunApp()
}
