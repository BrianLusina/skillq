package userv1

import "github.com/BrianLusina/skillq/server/app/internal/domain/ports/inbound"

// mapUserToUserResponse maps a user response to a user response dto
func mapUserToUserResponse(user inbound.UserResponse) userResponseDto {
	return userResponseDto{
		UUID:      user.UUID,
		KeyID:     user.KeyID,
		XID:       user.XID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		DeletedAt: user.DeletedAt,
		Name:      user.Name,
		Email:     user.Email,
		JobTitle:  user.JobTitle,
		Skills:    user.Skills,
		ImageUrl:  user.ImageUrl,
	}
}
