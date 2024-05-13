package userv1

import "github.com/gofiber/fiber/v2"

// RegisterHandlers registers all the handlers for the user v1 endpoint
func (api *UserV1Api) RegisterHandlers(r *fiber.App) {
	userApiGroup := r.Group("/api/v1/users")

	userApiGroup.Post("/", api.HandleCreateUser)
}
