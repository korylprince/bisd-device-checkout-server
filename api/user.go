package api

import (
	"fmt"

	auth "github.com/korylprince/go-ad-auth"
)

//AuthConfig holds configuration for connecting to an authentication source
type AuthConfig struct {
	ADConfig *auth.Config
	Group    string
}

//User represents an Active Directory User
type User struct {
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
}

//Authenticate authenticates the given username and password against the given config,
//returning user information if successful, nil if unsuccessful, or an error if one occurred.
func Authenticate(config *AuthConfig, username, password string) (*User, error) {
	status, attrs, err := auth.LoginWithAttrs(username, password, config.Group, config.ADConfig, []string{"displayName"})
	if err != nil {
		return nil, err
	}

	if !status {
		return nil, nil
	}

	if attrs == nil || len(attrs["displayName"]) != 1 || attrs["displayName"][0] == "" {
		return nil, fmt.Errorf("displayName doesn't exist for username: %s", username)
	}

	return &User{Username: username, DisplayName: attrs["displayName"][0]}, nil
}
