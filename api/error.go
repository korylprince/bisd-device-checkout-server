package api

import "fmt"

//Error wraps errors in the API
type Error struct {
	Description  string
	Err          error
	RequestError bool
}

func (e *Error) Error() string {
	if e.RequestError {
		if e.Err == nil {
			return fmt.Sprintf("Client Error: %s", e.Description)
		}
		return fmt.Sprintf("Client Error: %s: %v", e.Description, e.Err)
	}
	if e.Err == nil {
		return fmt.Sprintf("Server Error: %s", e.Description)
	}
	return fmt.Sprintf("Server Error: %s: %v", e.Description, e.Err)
}
