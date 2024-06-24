package mockuser

import (
	"fmt"
	"time"

	"github.com/BrianLusina/skillq/server/app/internal/domain/entities/user"
	"github.com/BrianLusina/skillq/server/domain/entity"
	"github.com/BrianLusina/skillq/server/domain/id"
	"github.com/go-faker/faker/v4"
)

func MockUser() (*user.User, error) {
	name := faker.FirstName()
	email := faker.Email()
	imageUrl := faker.URL()
	skill := faker.Word()
	jobTitle := faker.Word()
	password := faker.Password()
	createdAt := time.Now()
	updatedAt := time.Now()

	existingUser, err := user.New(user.UserParams{
		EntityParams: entity.EntityParams{
			EntityIDParams: entity.EntityIDParams{
				UUID:  id.NewUUID(),
				KeyID: id.NewKeyID(),
				XID:   id.NewXid(),
			},
			EntityTimestampParams: entity.EntityTimestampParams{
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
			},
			Metadata: map[string]any{},
		},
		Name:     name,
		Email:    email,
		ImageUrl: imageUrl,
		Skills:   []string{skill},
		JobTitle: jobTitle,
		Password: password,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create mock user: %v", err)
	}

	return &existingUser, nil
}
