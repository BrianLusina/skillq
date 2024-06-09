package userv1

import "github.com/gofiber/fiber/v2"

// RegisterHandlers registers all the handlers for the user v1 endpoint
func (api *UserV1Api) RegisterHandlers(app *fiber.App) {
	userApiGroup := app.Group("/api/v1/users")

	userApiGroup.Post("/", api.HandleCreateUser)
	userApiGroup.Post("/verify-email", api.HandleVerifyUserEmail)
	userApiGroup.Get("/:id", api.HandleGetUserById)
	userApiGroup.Get("/", api.HandleGetAllUsers)
	userApiGroup.Get("/skill/:skill", api.HandleGetAllUsersBySkill)
	userApiGroup.Delete("/:id", api.HandleDeleteUser)
}
