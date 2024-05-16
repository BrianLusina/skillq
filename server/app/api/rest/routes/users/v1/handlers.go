package userv1

import (
	"fmt"

	"github.com/BrianLusina/skillq/server/app/internal/domain/ports/inbound"
	"github.com/BrianLusina/skillq/server/app/internal/domain/ports/inbound/common"
	"github.com/BrianLusina/skillq/server/utils/tools"
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
		Name:     payload.Name,
		Email:    payload.Email,
		Password: payload.Password,
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

// HandleGetAllUsers gets all users
func (api *UserV1Api) HandleGetAllUsers(c *fiber.Ctx) error {
	ctx := c.Context()
	order := c.Query("order", string(common.CREATED_AT))
	sort := c.Query("sortby", string(common.DESC))
	limit := c.QueryInt("limit", 100)
	offset := c.QueryInt("offset", 0)

	params := common.NewRequestParams(
		common.WithRequestLimit(limit),
		common.WithOffset(offset),
		common.WithOrderBy(common.OrderBy(order)),
		common.WithSortOrder(common.SortOrder(sort)),
	)

	users, err := api.userService.GetAllUsers(ctx, params)
	if err != nil {
		// TODO: handle different types of error
		api.logger.Errorf("handler: failed to fetch users: %v", err)
		// utils.WriteWithError(w, http.StatusInternalServerError, "failed to create user")
		return err
	}

	response := tools.Map(users, func(u inbound.UserResponse, _ int) userResponseDto {
		return mapUserToUserResponse(u)
	})

	return c.JSON(response)
}

// HandleGetAllUsersBySkill gets all users with a given skill
func (api *UserV1Api) HandleGetAllUsersBySkill(c *fiber.Ctx) error {
	ctx := c.Context()

	skill := c.Params("skill")
	order := c.Query("order", string(common.CREATED_AT))
	sort := c.Query("sortby", string(common.DESC))
	limit := c.QueryInt("limit", 100)
	offset := c.QueryInt("offset", 0)

	params := common.NewRequestParams(
		common.WithRequestLimit(limit),
		common.WithOffset(offset),
		common.WithOrderBy(common.OrderBy(order)),
		common.WithSortOrder(common.SortOrder(sort)),
	)

	users, err := api.userService.GetAllUsersBySkill(ctx, skill, params)
	if err != nil {
		// TODO: handle different types of error
		api.logger.Errorf("handler: failed to fetch users by skill %s, err: %v", skill, err)
		// utils.WriteWithError(w, http.StatusInternalServerError, "failed to create user")
		return err
	}

	response := tools.Map(users, func(u inbound.UserResponse, _ int) userResponseDto {
		return mapUserToUserResponse(u)
	})

	return c.JSON(response)
}

// HandleGetAllUsersBySkill gets all users with a given skill
func (api *UserV1Api) HandleDeleteUser(c *fiber.Ctx) error {
	ctx := c.Context()

	userId := c.Params("id")

	err := api.userService.DeleteUser(ctx, userId)
	if err != nil {
		// TODO: handle different types of error
		api.logger.Errorf("handler: failed to delete user with ID %s, err: %v", userId, err)
		// utils.WriteWithError(w, http.StatusInternalServerError, "failed to create user")
		return err
	}

	return c.JSON(fiber.Map{
		"Message": fmt.Sprintf("Successfully deleted user %s", userId),
	})
}
