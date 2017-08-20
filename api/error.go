package api

import "fmt"

//Error wraps errors in the API
type Error struct {
	Description  string
	Err          error
	RequestError bool
}

func (e *Error) Error() string {
	return fmt.Sprintf("Server Error: %s: %v", e.Description, e.Err)
}
