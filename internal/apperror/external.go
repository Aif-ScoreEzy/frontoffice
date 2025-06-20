package apperror

import "fmt"

// ExternalAPIError is used to wrap error response from external APIs.
type ExternalAPIError struct {
	StatusCode int
	Message    string
}

func (e *ExternalAPIError) Error() string {
	return fmt.Sprintf("external api error [%d]: %s", e.StatusCode, e.Message)
}

// func (e *ExternalAPIError) IsUnauthorized() bool {
// 	return e.StatusCode == 401
// }

// func (e *ExternalAPIError) IsNotFound() bool {
// 	return e.StatusCode == 404
// }
