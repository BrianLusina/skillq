package v1

import (
	"encoding/json"
	"net/http"
)

// HandleCreateUser create a user
func (api *UserV1Api) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var request userRequestDto
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		// handle decode error
		return
	}
}
