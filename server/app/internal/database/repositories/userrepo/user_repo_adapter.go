package userrepo

import (
	"context"
	"fmt"
	"strings"

	"github.com/BrianLusina/skillq/server/app/internal/database/models"
	"github.com/BrianLusina/skillq/server/app/internal/domain/entities/user"
	"github.com/BrianLusina/skillq/server/app/internal/domain/ports/inbound/common"
	"github.com/BrianLusina/skillq/server/app/internal/domain/ports/outbound/repositories"
	"github.com/BrianLusina/skillq/server/domain/id"
	"github.com/BrianLusina/skillq/server/infra/mongodb"
	"github.com/BrianLusina/skillq/server/utils/tools"
	"github.com/pkg/errors"
)

// userRepoAdapter is the user repository adapter structure for managing user data
type userRepoAdapter struct {
	// dbClient is the database client used to handle connections to the database
	dbClient mongodb.MongoDBClient[models.UserModel]
}

var _ repositories.UserRepoPort = (*userRepoAdapter)(nil)

// New creates a new user repository adapter
func New(dbClient mongodb.MongoDBClient[models.UserModel]) repositories.UserRepoPort {
	return &userRepoAdapter{
		dbClient: dbClient,
	}
}

// CreateUser creates a user in the repository
func (repo *userRepoAdapter) CreateUser(ctx context.Context, userEntity user.User) (*user.User, error) {
	userModel := mapUserToModel(userEntity)
	_, err := repo.dbClient.Insert(ctx, userModel)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create user")
	}

	u, err := mapModelToUser(userModel)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// GetUserByUUID retrieves a user given their UUID
func (repo *userRepoAdapter) GetUserByUUID(ctx context.Context, userUUID id.UUID) (*user.User, error) {
	existingUser, err := repo.dbClient.FindById(ctx, "uuid", userUUID.String())
	if err != nil {
		return nil, errors.Wrapf(err, "failed to retrieve user by UUID %v", userUUID)
	}

	u, err := mapModelToUser(existingUser)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (repo *userRepoAdapter) GetAllUsers(ctx context.Context, params common.RequestParams) ([]user.User, error) {
	users, err := repo.dbClient.FindAll(ctx, mongodb.FilterOptions{
		Limit:     params.Limit,
		Offset:    params.Offset,
		SortOrder: mongodb.SortOrder(params.OrderOption.SortOrder),
		OrderBy:   string(params.OrderOption.OrderBy),
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to retrieve all users")
	}

	return tools.MapWithError(users, func(u models.UserModel, _ int) (user.User, error) {
		return mapModelToUser(u)
	})
}

func (repo *userRepoAdapter) GetAllUsersBySkill(ctx context.Context, skill string, params common.RequestParams) ([]user.User, error) {
	filterValues := map[string]map[string]string{
		"skills": {
			"$regex":   fmt.Sprintf("(?i)%s", skill),
			"$options": "i",
		},
	}

	users, err := repo.dbClient.FindAll(ctx, mongodb.FilterOptions{
		Limit:       params.Limit,
		Offset:      params.Offset,
		SortOrder:   mongodb.SortOrder(params.OrderOption.SortOrder),
		OrderBy:     string(params.OrderOption.OrderBy),
		FieldFilter: filterValues,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to retrieve all users by skill %s", skill)
	}

	usersWithSkill := tools.Filter(users, func(user models.UserModel) bool {
		hasSkill := false
		for _, userSkill := range user.Skills {
			if strings.ToLower(userSkill) == skill {
				hasSkill = true
				break
			}
			hasSkill = false
		}
		return hasSkill
	})

	return tools.MapWithError(usersWithSkill, func(u models.UserModel, _ int) (user.User, error) {
		return mapModelToUser(u)
	})
}

// UpdateUser updates a user given their ID
func (repo *userRepoAdapter) UpdateUser(ctx context.Context, userToUpdate user.User) (*user.User, error) {
	userModel := mapUserToModel(userToUpdate)

	err := repo.dbClient.Update(ctx, userModel, mongodb.UpdateOptions{
		Upsert: false,
		FieldOptions: map[string]any{
			"name":     userToUpdate.Name(),
			"email":    userToUpdate.Email(),
			"jobTitle": userToUpdate.JobTitle(),
			"imageUrl": userToUpdate.ImageUrl(),
			"skills":   userToUpdate.Skills(),
		},
		FilterParams: mongodb.FilterParams{
			Key:   "uuid",
			Value: userToUpdate.UUID().String(),
		},
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to update user %v", userToUpdate)
	}

	return &userToUpdate, nil
}

// DeleteUserById deletes a given user by the ID
func (repo *userRepoAdapter) DeleteUserById(ctx context.Context, userID id.UUID) error {
	err := repo.dbClient.Delete(ctx, "uuid", userID.String())
	if err != nil {
		return errors.Wrapf(err, "failed to delete user with ID: %s", userID)
	}
	return nil
}
