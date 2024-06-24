package usersvc

import (
	"github.com/BrianLusina/skillq/server/app/internal/domain/entities/user"
	"github.com/BrianLusina/skillq/server/app/internal/domain/ports/inbound"
)

func mapUserToUserResponse(userEntity user.User) *inbound.UserResponse {
	return &inbound.UserResponse{
		UUID:      userEntity.UUID().String(),
		KeyID:     userEntity.KeyID().String(),
		XID:       userEntity.XID().String(),
		CreatedAt: userEntity.CreatedAt(),
		UpdatedAt: userEntity.UpdatedAt(),
		DeletedAt: userEntity.DeletedAt(),
		Metadata:  userEntity.Metadata(),
		Name:      userEntity.Name(),
		Email:     userEntity.Email(),
		ImageUrl:  userEntity.ImageUrl(),
		Skills:    userEntity.Skills(),
		JobTitle:  userEntity.JobTitle(),
	}
}
