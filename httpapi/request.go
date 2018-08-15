package httpapi

//AuthenticateRequest is an username/password authentication request
type AuthenticateRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
