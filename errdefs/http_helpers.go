package errdefs

import (
	"net/http"
)

// GetHTTPErrorStatusCode retrieves status code from error message.
func GetHTTPErrorStatusCode(err error) int {
	if err == nil {
		return http.StatusInternalServerError
	}

	var statusCode int

	switch {
	case IsNotFound(err):
		statusCode = http.StatusNotFound
	case IsForbidden(err):
		statusCode = http.StatusForbidden
	}

	if statusCode == 0 {
		statusCode = http.StatusInternalServerError
	}

	return statusCode
}
