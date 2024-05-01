package test

import (
	"context"
	"fmt"
	"log"
	"testing"

	infra "github.com/BrianLusina/skillq/server/infra/mongodb"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	MONGODB_VERSION = "mongo:7.0.9"
	MONGODB_PORT    = "27017"
)

type TestConfig struct {
	Version    string
	Port       string
	Username   string
	Password   string
	Database   string
	Collection string
}

type TestDatabase struct {
}

// CreateMongoDbContainer creates a test mongo db container with the provided config
func CreateMongoDBContainer(ctx context.Context, testConfig TestConfig) (testcontainers.Container, error) {
	env := map[string]string{
		"MONGO_INITDB_DATABASE":      testConfig.Database,
		"MONGO_INITDB_ROOT_USERNAME": testConfig.Username,
		"MONGO_INITDB_ROOT_PASSWORD": testConfig.Password,
	}

	exposedPort := fmt.Sprintf("%s/tcp", testConfig.Port)

	req := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        testConfig.Version,
			ExposedPorts: []string{exposedPort},
			WaitingFor:   wait.ForListeningPort(nat.Port(testConfig.Port)),
			Env:          env,
		},
		Started: true,
	}

	container, err := testcontainers.GenericContainer(ctx, req)
	if err != nil {
		return &testcontainers.DockerContainer{}, fmt.Errorf("failed to create container: %v", err)
	}

	port, err := container.MappedPort(ctx, nat.Port(testConfig.Port))
	if err != nil {
		return &testcontainers.DockerContainer{}, fmt.Errorf("failed to retrieved mapped port: %v", err)
	}

	host, err := container.Host(ctx)
	if err != nil {
		return &testcontainers.DockerContainer{}, fmt.Errorf("failed to retrieve host: %v", err)
	}

	log.Printf("mongo db container ready and running on host: %s and port: %v ", host, port)

	return container, nil
}

func MongoDBClient[T any](t *testing.T, config infra.MongoDBConfig) infra.MongoDBClient[T] {
	client, err := infra.New[T](config)
	if err != nil {
		t.Errorf("failed to create mongodb client: %v", err)
		t.FailNow()
	}

	return client
}
