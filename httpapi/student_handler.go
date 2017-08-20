package httpapi

import (
	"net/http"

	"github.com/korylprince/bisd-device-checkout-server/api"
)

//GET /students
func handleReadStudentList(w http.ResponseWriter, r *http.Request) *handlerResponse {
	students, err := api.GetStudentList(r.Context())
	if resp := checkAPIError(err); resp != nil {
		return resp
	}

	return &handlerResponse{Code: http.StatusOK, Body: ReadStudentListResponse{Students: students}}
}
