package httpapi

//AuthenticateRequest is an username/password authentication request
type AuthenticateRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

//CheckoutRequest is a request to check out a device in the inventory
type CheckoutRequest struct {
	UserID string `json:"user_id"`
}
