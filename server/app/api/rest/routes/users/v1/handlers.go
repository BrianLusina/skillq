package userv1

import (
	"github.com/BrianLusina/skillq/server/app/internal/domain/ports/inbound"
	"github.com/gofiber/fiber/v2"
)

// HandleCreateUser create a user
func (api *UserV1Api) HandleCreateUser(c *fiber.Ctx) error {
	ctx := c.Context()

	payload := new(userRequestDto)
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

	response := mapUserToUserResponse(*user)

	return c.JSON(response)
}

// HandleGetUserById gets a user by an ID
func (api *UserV1Api) HandleGetUserById(c *fiber.Ctx) error {
	ctx := c.Context()
	userId := c.Params("id")

	user, err := api.userService.GetUserByUUID(ctx, userId)
	if err != nil {
		// TODO: handle different types of error
		api.logger.Errorf("handler: failed to fetch user: %v", err)
		// utils.WriteWithError(w, http.StatusInternalServerError, "failed to create user")
		return err
	}

	response := mapUserToUserResponse(*user)

	return c.JSON(response)
}
