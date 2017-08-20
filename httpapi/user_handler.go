package httpapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/korylprince/bisd-device-checkout-server/api"
)

//POST /auth
func handleAuthenticate(config *api.AuthConfig, s SessionStore) returnHandler {
	return func(w http.ResponseWriter, r *http.Request) *handlerResponse {
		var req *AuthenticateRequest
		d := json.NewDecoder(r.Body)

		err := d.Decode(&req)
		if err != nil || req == nil {
			return handleError(http.StatusBadRequest, fmt.Errorf("Could not decode json: %v", err))
		}

		if req.Username == "" || req.Password == "" {
			return handleError(http.StatusBadRequest, errors.New("username or password empty"))
		}

		user, err := api.Authenticate(config, req.Username, req.Password)
		if err != nil {
			return handleError(http.StatusUnauthorized, fmt.Errorf("Could not authenticate user %s: %v", req.Username, err))
		}
		if user == nil {
			return handleError(http.StatusUnauthorized, errors.New("Bad username or password"))
		}

		key, err := s.Create(user)
		if err != nil {
			return handleError(http.StatusInternalServerError, fmt.Errorf("Could not create session: %v", err))
		}

		return &handlerResponse{Code: http.StatusOK, Body: &AuthenticateResponse{SessionKey: key, User: user}}
	}
}
