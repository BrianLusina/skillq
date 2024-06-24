package minio

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/BrianLusina/skillq/server/infra/logger"
	"github.com/BrianLusina/skillq/server/infra/storage"
	miniotest "github.com/BrianLusina/skillq/server/infra/storage/minio/test"
	"github.com/stretchr/testify/assert"
)

func TestMinioStorageClient(t *testing.T) {
	ctx := context.Background()
	minioTestContainerConfig := miniotest.TestContainerConfig{
		Version:     miniotest.MINIO_TEST_CONTAINER_TAG,
		AccessKeyID: miniotest.MINIO_TEST_ACCESS_KEY,
		SecretKey:   miniotest.MINIO_TEST_SECRET_KEY,
	}

	container, err := miniotest.MinioContainer(ctx, minioTestContainerConfig)
	if err != nil {
		t.Fatalf("failed to start minio test container")
		t.FailNow()
	}

	defer func() {
		if err := container.Terminate(ctx); err != nil {
			log.Fatalf("failed to terminate container: %s", err)
		}
	}()

	endpoint, err := container.Host(ctx)
	if err != nil {
		t.Fatalf("failed to retrieve minio host")
		t.FailNow()
	}

	ports, err := container.Ports(ctx)
	if err != nil {
		t.Fatalf("failed to retrieve ports")
		t.FailNow()
	}
	firstPort := "9000"
	for _, port := range ports {
		for _, p := range port {
			firstPort = p.HostPort
			break
		}
	}

	config := Config{
		Endpoint:        fmt.Sprintf("%s:%s", endpoint, firstPort),
		AccessKeyID:     container.Username,
		SecretAccessKey: container.Password,
	}
	log := logger.New()

	minioClient, err := NewClient(config, log)
	assert.NoError(t, err)
	assert.NotNil(t, minioClient)

	// tests
	t.Run("should upload a file for storage", func(t *testing.T) {
		storageItem := storage.StorageItem{
			Name:        "document-name",
			Content:     "data:text/plain;charset=utf-8;base64,aGV5YQ==",
			ContentType: "text/plain",
			Bucket:      "images",
			Metadata: map[string]string{
				"profile": "profile photo",
				"size":    "64MB",
			},
		}

		// TODO: check location url is returned. logs indicate that the item is stored, but the location is not returned
		_, err := minioClient.Upload(ctx, storageItem)
		assert.NoError(t, err)
	})
}
