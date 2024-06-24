package mockuser

import (
	"time"

	"github.com/BrianLusina/skillq/server/app/internal/domain/entities/user"
	"github.com/BrianLusina/skillq/server/domain/id"
)

func MockUserVerification(code string) user.UserVerification {
	createdAt := time.Now()
	updatedAt := time.Now()

	verification := user.NewVerification(user.UserVerificationParams{
		ID:         id.NewUUID(),
		Code:       code,
		UserId:     id.NewUUID(),
		IsVerified: false,
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
	})

	return verification
}
