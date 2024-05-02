package v1

import (
	"encoding/json"
	"net/http"

	"github.com/BrianLusina/skillq/server/app/api/rest/utils"
	"github.com/BrianLusina/skillq/server/app/internal/domain/ports/inbound"

	"github.com/BrianLusina/skillq/server/infra/logger"
)

// HandleCreateUser create a user
func (api *UserV1Api) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger := logger.FromContext(ctx)

	var payload userRequestDto
	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		logger.Errorf("userapi/v1 handler: failed to decode request: %v", err)
		utils.HandleDecodeErr(w, err)
		return
	}

	defer r.Body.Close()

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
		logger.Errorf("handler: failed to create user: %v", err)
		utils.WriteWithError(w, http.StatusInternalServerError, "failed to create user")
		return
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

	utils.WriteWithStatus(w, http.StatusCreated, response)
}
