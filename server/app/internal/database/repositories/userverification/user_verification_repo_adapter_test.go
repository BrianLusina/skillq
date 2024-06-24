package userverificationrepo

import (
	"context"
	"errors"
	"testing"

	"github.com/BrianLusina/skillq/server/app/internal/database/models"
	mockuser "github.com/BrianLusina/skillq/server/app/internal/domain/entities/user/mocks"
	mockmongodb "github.com/BrianLusina/skillq/server/infra/mongodb/mocks"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/mock/gomock"
)

func TestUserVerificationRepoAdapter(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockDbClient := mockmongodb.NewMockMongoDBClient[models.UserVerificationModel](mockCtrl)
	userVerificationRepositoryAdapter := userVerificationRepoAdapter{dbClient: mockDbClient}
	adapter := New(mockDbClient)
	assert.NotNil(t, adapter)

	ctx := context.Background()

	t.Run("creating a user verification", func(t *testing.T) {
		t.Run("should return error when there is a failure to create user verication", func(t *testing.T) {
			defer mockCtrl.Finish()

			dbErr := errors.New("failed to create user verification")

			mockDbClient.EXPECT().Insert(ctx, gomock.Any()).Return(primitive.ObjectID{}, dbErr).Times(1)

			testUserVerification := mockuser.MockUserVerification("some_code")

			actual, err := userVerificationRepositoryAdapter.CreateUserVerification(ctx, testUserVerification)
			assert.Error(t, err)
			assert.Nil(t, actual)
		})

		t.Run("should return created user verification when there is a success in creating user verification", func(t *testing.T) {
			defer mockCtrl.Finish()

			mockDbClient.EXPECT().Insert(ctx, gomock.Any()).Return(primitive.ObjectID{}, nil).Times(1)

			testUserVerification := mockuser.MockUserVerification("test_code")

			actual, err := userVerificationRepositoryAdapter.CreateUserVerification(ctx, testUserVerification)
			assert.NoError(t, err)
			assert.NotNil(t, actual)
		})
	})

	t.Run("get user verification", func(t *testing.T) {
		testUserVerification := mockuser.MockUserVerification("somecode")

		testUserVerificationModel := mapUserVerificationToModel(testUserVerification)

		t.Run("by ID", func(t *testing.T) {

			t.Run("should return nil & error when there is a failure to retrieve user verification by ID", func(t *testing.T) {
				defer mockCtrl.Finish()

				userVerificationID := testUserVerification.ID()

				dbError := errors.New("failed to retrieve user verification")
				mockDbClient.EXPECT().FindById(ctx, "uuid", userVerificationID.String()).Return(models.UserVerificationModel{}, dbError).Times(1)

				actualUser, err := userVerificationRepositoryAdapter.GetUserVerificationByUUID(ctx, userVerificationID)
				assert.Error(t, err)
				assert.Nil(t, actualUser)
			})

			t.Run("should return user verification model & nill error when there is a success in retrieving a user verification by UUID", func(t *testing.T) {
				defer mockCtrl.Finish()

				userVerificationID := testUserVerification.ID()

				mockDbClient.EXPECT().FindById(ctx, "uuid", userVerificationID.String()).Return(testUserVerificationModel, nil).Times(1)

				actualUser, err := userVerificationRepositoryAdapter.GetUserVerificationByUUID(ctx, userVerificationID)
				assert.NoError(t, err)
				assert.NotNil(t, actualUser)
			})
		})

		t.Run("by a code", func(t *testing.T) {
			t.Run("should return nil & error when there is a failure to retrieve verification for a given code", func(t *testing.T) {
				defer mockCtrl.Finish()

				code := "somecode"

				dbError := errors.New("failed to retrieve user verification")
				mockDbClient.EXPECT().FindById(ctx, "code", code).Return(models.UserVerificationModel{}, dbError).Times(1)

				actualUserVerification, err := userVerificationRepositoryAdapter.GetUserVerificationByCode(ctx, code)
				assert.Error(t, err)
				assert.Nil(t, actualUserVerification)
			})

			t.Run("should return user verification & nil error when there is a success in retrieving user verification", func(t *testing.T) {
				defer mockCtrl.Finish()

				code := "somecode"

				mockDbClient.EXPECT().FindById(ctx, "code", code).Return(testUserVerificationModel, nil).Times(1)

				actualUsers, err := userVerificationRepositoryAdapter.GetUserVerificationByCode(ctx, code)
				assert.NoError(t, err)
				assert.NotNil(t, actualUsers)
			})
		})
	})

}
