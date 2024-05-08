package userv1

import (
	"github.com/BrianLusina/skillq/server/app/internal/domain/ports/inbound"
	"github.com/gofiber/fiber/v2"
)

// HandleCreateUser create a user
func (api *UserV1Api) HandleCreateUser(c *fiber.Ctx) error {
	ctx := c.Context()

	payload := new(userRequestDto)
	c.BodyParser(payload)
	if err := c.BodyParser(payload); err != nil {
		api.logger.Errorf("userapi/v1 handler: failed to decode request: %v", err)
		// TODO: format error to DTO
		return err
	}

	userRequest := inbound.UserRequest{
		Name:  payload.Name,
		Email: payload.Email,
		Image: inbound.UserImageRequest{
			Type:    payload.Image.ImageType,
			Content: payload.Image.Content,
		},
		Skills:   payload.Skills,
		JobTitle: payload.JobTitle,
	}
	user, err := api.userService.CreateUser(ctx, userRequest)
	if err != nil {
		// TODO: handle different types of error
		api.logger.Errorf("handler: failed to create user: %v", err)
		// utils.WriteWithError(w, http.StatusInternalServerError, "failed to create user")
		return err
	}

	response := userResponseDto{
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

	return c.JSON(response)
}
