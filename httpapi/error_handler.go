package httpapi

import (
	"errors"
	"net/http"

	"github.com/korylprince/bisd-device-checkout-server/api"
)

// ErrorResponse represents an HTTP error
type ErrorResponse struct {
	Code        int    `json:"code"`
	Error       string `json:"error"`
	Description string `json:"description"`
}

// handleError returns a handlerResponse response for the given code
func handleError(code int, err error) *handlerResponse {
	return &handlerResponse{Code: code, Body: &ErrorResponse{Code: code, Error: http.StatusText(code), Description: err.Error()}, Err: err}
}

// notFoundHandler returns a 401 handlerResponse
func notFoundHandler(_ http.ResponseWriter, _ *http.Request) *handlerResponse {
	return handleError(http.StatusNotFound, errors.New("Could not find handler"))
}

// checkAPIError checks an api.Error and returns a handlerResponse for it, or nil if there was no error
func checkAPIError(err error) *handlerResponse {
	if err == nil {
		return nil
	}

	if e, ok := err.(*api.Error); ok && e.RequestError {
		return handleError(http.StatusBadRequest, err)
	}
	return handleError(http.StatusInternalServerError, err)
}
