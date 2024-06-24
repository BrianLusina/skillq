package userrepo

import (
	"context"
	"log"
	"testing"

	"github.com/BrianLusina/skillq/server/app/internal/database/models"
	mockuser "github.com/BrianLusina/skillq/server/app/internal/domain/entities/user/mocks"
	"github.com/BrianLusina/skillq/server/infra/mongodb"
	"github.com/BrianLusina/skillq/server/infra/mongodb/test"
	"github.com/stretchr/testify/assert"
)

func TestUserRepoAdapterIntegration(t *testing.T) {
	ctx := context.Background()
	testConfig := test.TestConfig{
		Version:    test.MONGODB_VERSION,
		Port:       test.MONGODB_PORT,
		Username:   "skillq-user",
		Password:   "skillq-password",
		Database:   "skillq",
		Collection: "users",
	}
	mongodbContainer, err := test.CreateMongoDBContainer(ctx, testConfig)
	if err != nil {
		t.Errorf("failed to start mongodb container: %s", err)
		t.FailNow()
	}

	// Clean up the container
	defer func() {
		if err := mongodbContainer.Terminate(ctx); err != nil {
			log.Fatalf("failed to terminate mongodb container: %s", err)
		}
	}()

	host, err := mongodbContainer.Host(ctx)
	if err != nil {
		t.Errorf("failed to retrieve host from mongodb container: %v", err)
		t.FailNow()
	}

	config := mongodb.MongoDBConfig{
		Client: mongodb.ClientOptions{
			Host:     host,
			User:     testConfig.Username,
			Password: testConfig.Password,
		},
		DBConfig: mongodb.DatabaseConfig{
			DatabaseName:   testConfig.Database,
			CollectionName: testConfig.Collection,
		},
	}

	mongoDbClient := test.MongoDBClient[models.UserModel](t, config)
	userRepositoryAdapter := userRepoAdapter{
		dbClient: mongoDbClient,
	}

	t.Run("creating a user", func(t *testing.T) {
		t.Run("should return created user when there is a success in creating user", func(t *testing.T) {
			defer mongoDbClient.Disconnect(ctx)

			testUser, err := mockuser.MockUser()
			assert.NoError(t, err)

			actualUser, actualErr := userRepositoryAdapter.CreateUser(ctx, *testUser)
			assert.NoError(t, actualErr)
			assert.NotNil(t, actualUser)
		})
	})

}
