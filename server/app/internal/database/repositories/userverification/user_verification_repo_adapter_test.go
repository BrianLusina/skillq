package userverificationrepo

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/BrianLusina/skillq/server/app/internal/database/models"
	"github.com/BrianLusina/skillq/server/app/internal/domain/entities/user"
	mockuser "github.com/BrianLusina/skillq/server/app/internal/domain/entities/user/mocks"
	"github.com/BrianLusina/skillq/server/app/internal/domain/ports/inbound/common"
	mockmongodb "github.com/BrianLusina/skillq/server/infra/mongodb/mocks"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/mock/gomock"
)

func TestUserVerificationRepoAdapter(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockDbClient := mockmongodb.NewMockMongoDBClient[models.UserVerificationModel](mockCtrl)
	userRepositoryAdapter := userVerificationRepoAdapter{dbClient: mockDbClient}
	adapter := New(mockDbClient)
	assert.NotNil(t, adapter)

	ctx := context.Background()

	t.Run("creating a user verification", func(t *testing.T) {
		t.Run("should return error when there is a failure to create user verication", func(t *testing.T) {
			defer mockCtrl.Finish()

			dbErr := errors.New("failed to create user verification")

			mockDbClient.EXPECT().Insert(ctx, gomock.Any()).Return(primitive.ObjectID{}, dbErr).Times(1)

			testUserVerification := mockuser.MockUserVerification("some_code")

			actual, err := userRepositoryAdapter.CreateUserVerification(ctx, testUserVerification)
			assert.Error(t, err)
			assert.Nil(t, actual)
		})

		t.Run("should return created user verification when there is a success in creating user verification", func(t *testing.T) {
			defer mockCtrl.Finish()

			mockDbClient.EXPECT().Insert(ctx, gomock.Any()).Return(primitive.ObjectID{}, nil).Times(1)

			testUserVerification := mockuser.MockUserVerification("test_code")

			actual, err := userRepositoryAdapter.CreateUserVerification(ctx, testUserVerification)
			assert.NoError(t, err)
			assert.NotNil(t, actual)
		})
	})

	t.Run("Get user", func(t *testing.T) {
		testUserVerification := mockuser.MockUserVerification("somecode")

		testUserVerificationModel := mapUserVerificationToModel(testUserVerification)

		t.Run("by ID", func(t *testing.T) {

			t.Run("should return nil & error when there is a failure to retrieve user verification by ID", func(t *testing.T) {
				defer mockCtrl.Finish()

				userVerificationID := testUserVerification.ID()

				dbError := errors.New("failed to retrieve user verification")
				mockDbClient.EXPECT().FindById(ctx, "uuid", userVerificationID.String()).Return(models.UserVerificationModel{}, dbError).Times(1)

				actualUser, err := userRepositoryAdapter.GetUserVerificationByUUID(ctx, userVerificationID)
				assert.Error(t, err)
				assert.Nil(t, actualUser)
			})

			t.Run("should return user model & error when there is a success in retrieving a user by UUID", func(t *testing.T) {
				defer mockCtrl.Finish()

				userUUID := testUserVerification.UUID()

				mockDbClient.EXPECT().FindById(ctx, "uuid", userUUID.String()).Return(testUserVerificationModel, nil).Times(1)

				actualUser, err := userRepositoryAdapter.GetUserByUUID(ctx, userUUID)
				assert.NoError(t, err)
				assert.NotNil(t, actualUser)
			})
		})
	})

	t.Run("Get all users", func(t *testing.T) {
		testUserOne, err := mockuser.MockUser()
		assert.NoError(t, err)
		testUserOneModel := mapUserToModel(*testUserOne)

		testUserTwo, err := mockuser.MockUser()
		assert.NoError(t, err)
		testUserTwoModel := mapUserToModel(*testUserTwo)

		testUserThree, err := mockuser.MockUser()
		assert.NoError(t, err)
		testUserThreeModel := mapUserToModel(*testUserThree)

		testUsers := []user.User{*testUserOne, *testUserTwo, *testUserThree}
		testUserModels := []models.UserModel{testUserOneModel, testUserTwoModel, testUserThreeModel}

		t.Run("should return nil & error when there is a failure to retrieve all users", func(t *testing.T) {
			defer mockCtrl.Finish()

			dbError := errors.New("failed to retrieve users")
			mockDbClient.EXPECT().FindAll(ctx, map[string]map[string]string{}).Return(nil, dbError).Times(1)

			actualUsers, err := userRepositoryAdapter.GetAllUsers(ctx, common.NewRequestParams())
			assert.Error(t, err)
			assert.Nil(t, actualUsers)
		})

		t.Run("should return users & nil error when there is a success in retrieving a users", func(t *testing.T) {
			defer mockCtrl.Finish()

			mockDbClient.EXPECT().FindAll(ctx, map[string]map[string]string{}).Return(testUserModels, nil).Times(1)

			actualUsers, err := userRepositoryAdapter.GetAllUsers(ctx, common.NewRequestParams())
			assert.NoError(t, err)
			assert.NotNil(t, actualUsers)
			assert.ElementsMatch(t, testUsers, actualUsers)
		})

		t.Run("by a skill", func(t *testing.T) {
			t.Run("should return nil & error when there is a failure to retrieve all users for a given skill", func(t *testing.T) {
				defer mockCtrl.Finish()

				skill := "hunter"
				filterValues := map[string]map[string]string{
					"skills": {
						"$regex":   fmt.Sprintf("(?i)%s", skill),
						"$options": "i",
					},
				}

				dbError := errors.New("failed to retrieve users")
				mockDbClient.EXPECT().FindAll(ctx, filterValues).Return(nil, dbError).Times(1)

				actualUsers, err := userRepositoryAdapter.GetAllUsersBySkill(ctx, skill, common.NewRequestParams())
				assert.Error(t, err)
				assert.Nil(t, actualUsers)
			})

			t.Run("should return users & nil error when there is a success in retrieving users", func(t *testing.T) {
				defer mockCtrl.Finish()

				skill := "hunter"
				filterValues := map[string]map[string]string{
					"skills": {
						"$regex":   fmt.Sprintf("(?i)%s", skill),
						"$options": "i",
					},
				}

				mockDbClient.EXPECT().FindAll(ctx, filterValues).Return(testUserModels, nil).Times(1)

				actualUsers, err := userRepositoryAdapter.GetAllUsersBySkill(ctx, skill, common.NewRequestParams())
				assert.NoError(t, err)
				assert.NotNil(t, actualUsers)
			})
		})
	})

}
