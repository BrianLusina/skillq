package userverificationrepo

import (
	"context"

	"github.com/BrianLusina/skillq/server/app/internal/database/models"
	"github.com/BrianLusina/skillq/server/app/internal/domain/entities/user"
	"github.com/BrianLusina/skillq/server/app/internal/domain/ports/outbound/repositories"
	"github.com/BrianLusina/skillq/server/domain/id"
	"github.com/BrianLusina/skillq/server/infra/mongodb"
	"github.com/pkg/errors"
)

// userVerificationRepoAdapter is the user verification repository adapter structure for managing user verification data
type userVerificationRepoAdapter struct {
	// dbClient is the database client used to handle connections to the database
	dbClient mongodb.MongoDBClient[models.UserVerificationModel]
}

var _ repositories.UserVerificationRepoPort = (*userVerificationRepoAdapter)(nil)

// New creates a new user verification repository adapter
func New(dbClient mongodb.MongoDBClient[models.UserVerificationModel]) repositories.UserVerificationRepoPort {
	return &userVerificationRepoAdapter{
		dbClient: dbClient,
	}
}

// CreateUserVerification creates a user verification in the repository
func (repo *userVerificationRepoAdapter) CreateUserVerification(ctx context.Context, userVerification user.UserVerification) (*user.UserVerification, error) {
	userVerificationModel := mapUserVerificationToModel(userVerification)
	_, err := repo.dbClient.Insert(ctx, userVerificationModel)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create user")
	}

	u, err := mapUserVerificationModelToEntity(userVerificationModel)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// GetUserVerificationByUUID retrieves a user verification given the code
func (repo *userVerificationRepoAdapter) GetUserVerificationByUUID(ctx context.Context, userVerificationUUID id.UUID) (*user.UserVerification, error) {
	existingUser, err := repo.dbClient.FindById(ctx, "uuid", userVerificationUUID.String())
	if err != nil {
		return nil, errors.Wrapf(err, "failed to retrieve user verification by UUID %v", userVerificationUUID)
	}

	u, err := mapUserVerificationModelToEntity(existingUser)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

// GetUserVerificationByCode retrieves a user verification given the code
func (repo *userVerificationRepoAdapter) GetUserVerificationByCode(ctx context.Context, code string) (*user.UserVerification, error) {
	userVerificationModel, err := repo.dbClient.FindById(ctx, "code", code)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to retrieve user verification by code %v", code)
	}

	u, err := mapUserVerificationModelToEntity(userVerificationModel)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

// UpdateUserVerification updates the user verification
func (repo *userVerificationRepoAdapter) UpdateUserVerification(ctx context.Context, request repositories.UpdateUserVerificationRequest) error {
	userVerification, err := repo.dbClient.FindById(ctx, "user_id", request.UserID.String())
	if err != nil {
		return errors.Wrapf(err, "failed to retrieve user verification for user id %s", request.UserID)
	}

	userVerification.IsVerified = request.IsVerified

	err = repo.dbClient.Update(ctx, userVerification, mongodb.UpdateOptions{
		Upsert: false,
		FieldOptions: map[string]any{
			"is_verified": request.IsVerified,
		},
		FilterParams: mongodb.FilterParams{
			Key:   "uuid",
			Value: userVerification.UUID,
		},
	})

	if err != nil {
		return errors.Wrapf(err, "failed to update user verification for user %s", request.UserID)
	}
	return nil
}