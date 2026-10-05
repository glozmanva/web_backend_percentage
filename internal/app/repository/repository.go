package repository

import (
	"context"
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Repository struct {
	db              *gorm.DB
	minio           *minio.Client
	minioBucketName string
	minioEndpoint   string
}

type RepositorySettings struct {
	PostgresDSN     string
	MinioEndpoint   string
	MinioAccessKey  string
	MinioSecretKey  string
	MinioBucketName string
}

func New(settings *RepositorySettings) (*Repository, error) {
	db, err := gorm.Open(
		postgres.Open(settings.PostgresDSN),
		&gorm.Config{},
	)
	if err != nil {
		return nil, err
	}

	minioClient, err := minio.New(
		settings.MinioEndpoint,
		&minio.Options{
			Creds: credentials.NewStaticV4(
				settings.MinioAccessKey,
				settings.MinioSecretKey,
				"",
			),
			Secure: false,
		},
	)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()

	exists, err := minioClient.BucketExists(
		ctx,
		settings.MinioBucketName,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"не удалось проверить bucket MinIO: %w",
			err,
		)
	}

	if !exists {
		err = minioClient.MakeBucket(
			ctx,
			settings.MinioBucketName,
			minio.MakeBucketOptions{},
		)
		if err != nil {
			return nil, fmt.Errorf(
				"не удалось создать bucket MinIO: %w",
				err,
			)
		}
	}

	return &Repository{
		db:              db,
		minio:           minioClient,
		minioBucketName: settings.MinioBucketName,
		minioEndpoint:   settings.MinioEndpoint,
	}, nil
}
