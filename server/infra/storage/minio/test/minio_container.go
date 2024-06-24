package miniotest

import (
	"context"

	"github.com/pkg/errors"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/minio"
)

const (
	MINIO_TEST_CONTAINER_TAG = "minio/minio:RELEASE.2024-05-01T01-11-10Z"
	MINIO_TEST_ACCESS_KEY    = "AKIAIOSFODNN7EXAMPLE"
	MINIO_TEST_SECRET_KEY    = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
	MINIO_PORT               = "9000"
)

type TestContainerConfig struct {
	Version     string
	Username    string
	Password    string
	AccessKeyID string
	SecretKey   string
	Port        string
}

// MinioContainer creates a minio test container for use in tests
func MinioContainer(ctx context.Context, config TestContainerConfig) (*minio.MinioContainer, error) {
	env := map[string]string{
		"MINIO_ACCESS_KEY": config.AccessKeyID,
		"MINIO_SECRET_KEY": config.SecretKey,
	}

	container, err := minio.RunContainer(ctx,
		testcontainers.WithImage(config.Version),
		testcontainers.WithEnv(env),
	)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create minio test container")
	}

	return container, nil
}
