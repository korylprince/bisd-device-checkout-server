package api

import (
	"fmt"

	auth "github.com/korylprince/go-ad-auth/v3"
)

// AuthConfig holds configuration for connecting to an authentication source
type AuthConfig struct {
	ADConfig *auth.Config
	Group    string
}

// User represents an Active Directory User
type User struct {
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
}

// Authenticate authenticates the given username and password against the given config,
// returning user information if successful, nil if unsuccessful, or an error if one occurred.
func Authenticate(config *AuthConfig, username, password string) (*User, error) {
	status, entry, groups, err := auth.AuthenticateExtended(config.ADConfig, username, password, []string{"displayName"}, []string{config.Group})
	if err != nil {
		return nil, fmt.Errorf("Error attempting to authenticate as %s: %v", username, err)
	}

	if !status {
		return nil, nil
	}

	if len(groups) == 0 {
		return nil, nil
	}

	if entry.GetAttributeValue("displayName") == "" {
		return nil, fmt.Errorf("displayName doesn't exist for username: %s", username)
	}

	return &User{Username: username, DisplayName: entry.GetAttributeValue("displayName")}, nil
}
