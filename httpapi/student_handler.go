package httpapi

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/korylprince/bisd-device-checkout-server/api"
)

//GET /students
func handleReadStudentList(w http.ResponseWriter, r *http.Request) *handlerResponse {
	type student struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		OtherID   string `json:"other_id"`
		Grade     int    `json:"grade"`
	}

	type response []*student

	students, err := api.GetStudentList(r.Context())
	if resp := checkAPIError(err); resp != nil {
		return resp
	}

	var list response
	for _, s := range students {
		list = append(list, &student{FirstName: s.FirstName, LastName: s.LastName, OtherID: s.OtherID, Grade: s.Grade})
	}

	return &handlerResponse{Code: http.StatusOK, Body: list}
}

//GET /students/:otherID/status
func handleReadStudentStatus(w http.ResponseWriter, r *http.Request) *handlerResponse {
	otherID := mux.Vars(r)["otherID"]

	student, err := api.GetStudent(r.Context(), otherID)
	if resp := checkAPIError(err); resp != nil {
		return resp
	}

	status, err := student.Status(r.Context())
	if resp := checkAPIError(err); resp != nil {
		return resp
	}

	return &handlerResponse{Code: http.StatusOK, Body: status}
}
