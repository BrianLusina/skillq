package userrepo

import (
	"github.com/BrianLusina/skillq/server/app/internal/database/models"
	"github.com/BrianLusina/skillq/server/app/internal/domain/entities/user"
	"github.com/BrianLusina/skillq/server/domain/entity"
	"github.com/BrianLusina/skillq/server/domain/id"
)

// mapUserToModel maps a user entity to a user model
func mapUserToModel(userEntity user.User) models.UserModel {
	return models.UserModel{
		BaseModel: models.BaseModel{
			UUID:      userEntity.UUID().String(),
			KeyID:     userEntity.KeyID().String(),
			XID:       userEntity.XID().String(),
			Metadata:  userEntity.Metadata(),
			CreatedAt: userEntity.CreatedAt(),
			UpdatedAt: userEntity.UpdatedAt(),
			DeletedAt: userEntity.DeletedAt(),
		},
		Name:         userEntity.Name(),
		Email:        userEntity.Email(),
		ImageUrl:     userEntity.ImageUrl(),
		JobTitle:     userEntity.JobTitle(),
		Skills:       userEntity.Skills(),
		PasswordHash: userEntity.Password(),
	}
}

// mapModelToUser maps a user model to a user entity
func mapModelToUser(userModel models.UserModel) (user.User, error) {
	keyId, err := id.StringToKeyID(userModel.BaseModel.KeyID)
	if err != nil {
		return user.User{}, err
	}

	uuid, err := id.StringToUUID(userModel.BaseModel.UUID)
	if err != nil {
		return user.User{}, err
	}

	xid, err := id.StringToXid(userModel.BaseModel.XID)
	if err != nil {
		return user.User{}, err
	}

	return user.New(user.UserParams{
		EntityParams: entity.EntityParams{
			EntityIDParams: entity.EntityIDParams{
				UUID:  uuid,
				KeyID: keyId,
				XID:   xid,
			},
			EntityTimestampParams: entity.EntityTimestampParams{
				CreatedAt: userModel.BaseModel.CreatedAt,
				UpdatedAt: userModel.BaseModel.UpdatedAt,
				DeletedAt: userModel.BaseModel.DeletedAt,
			},
			Metadata: userModel.BaseModel.Metadata,
		},
		Name:     userModel.Name,
		Email:    userModel.Email,
		Password: userModel.PasswordHash,
		JobTitle: userModel.JobTitle,
		Skills:   userModel.Skills,
		ImageUrl: userModel.ImageUrl,
	})
}

// mapUserVerificationToModel maps a user entity to a user model
func mapUserVerificationToModel(userVerificationEntity user.UserVerification) models.UserVerificationModel {
	return models.UserVerificationModel{
		UUID:       userVerificationEntity.ID().String(),
		UserId:     userVerificationEntity.UserID().String(),
		CreatedAt:  userVerificationEntity.CreatedAt(),
		UpdatedAt:  userVerificationEntity.UpdatedAt(),
		Code:       userVerificationEntity.Code(),
		IsVerified: userVerificationEntity.IsVerified(),
	}
}

// mapUserToModel maps a user entity to a user model
func mapUserVerificationModelToEntity(userVerificationModel models.UserVerificationModel) (user.UserVerification, error) {
	uuid, err := id.StringToUUID(userVerificationModel.UUID)
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
		CreatedAt:  userVerificationModel.CreatedAt,
		UpdatedAt:  userVerificationModel.UpdatedAt,
		Code:       userVerificationModel.Code,
		IsVerified: userVerificationModel.IsVerified,
	}), nil
}
