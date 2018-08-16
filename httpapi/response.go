package httpapi

import "github.com/korylprince/bisd-device-checkout-server/api"

//AuthenticateResponse is a successful authentication response including the session key and User
type AuthenticateResponse struct {
	SessionID string    `json:"session_id"`
	User      *api.User `json:"user"`
}

//ReadStudentListResponse is a response with a list of students
type ReadStudentListResponse struct {
	Students []*api.Student `json:"students"`
}
