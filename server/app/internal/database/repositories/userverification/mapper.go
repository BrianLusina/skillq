package userverificationrepo

import (
	"github.com/BrianLusina/skillq/server/app/internal/database/models"
	"github.com/BrianLusina/skillq/server/app/internal/domain/entities/user"
	"github.com/BrianLusina/skillq/server/domain/id"
)

// mapUserVerificationToModel maps a user entity to a user model
func mapUserVerificationToModel(userVerificationEntity user.UserVerification) models.UserVerificationModel {
	return models.UserVerificationModel{
		BaseModel: models.BaseModel{
			UUID:      userVerificationEntity.ID().String(),
			CreatedAt: userVerificationEntity.CreatedAt(),
			UpdatedAt: userVerificationEntity.UpdatedAt(),
		},
		UserId:     userVerificationEntity.UserID().String(),
		Code:       userVerificationEntity.Code(),
		IsVerified: userVerificationEntity.IsVerified(),
	}
}

// mapUserToModel maps a user entity to a user model
func mapUserVerificationModelToEntity(userVerificationModel models.UserVerificationModel) (user.UserVerification, error) {
	uuid, err := id.StringToUUID(userVerificationModel.BaseModel.UUID)
	if err != nil {
		return user.UserVerification{}, err
	}

	userId, err := id.StringToUUID(userVerificationModel.UserId)
	if err != nil {
		return user.UserVerification{}, err
	}

	return user.NewVerification(user.UserVerificationParams{
		ID:         uuid,
		UserId:     userId,
		CreatedAt:  userVerificationModel.BaseModel.CreatedAt,
		UpdatedAt:  userVerificationModel.BaseModel.UpdatedAt,
		Code:       userVerificationModel.Code,
		IsVerified: userVerificationModel.IsVerified,
	}), nil
}
