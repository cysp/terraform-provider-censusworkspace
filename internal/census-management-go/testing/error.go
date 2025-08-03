package testing

import (
	"fmt"
	"net/http"

	cm "github.com/cysp/terraform-provider-censusworkspace/internal/census-management-go"
)

var (
	errNotFound = newStatusMessageError(http.StatusNotFound, cm.ResponseStatusNotFound, "not found")
)

type statusMessageError struct {
	StatusCode int
	Response   cm.StatusResponse
}

var _ error = (*statusMessageError)(nil)

func newStatusMessageError(statusCode int, status cm.ResponseStatus, message string) statusMessageError {
	return statusMessageError{
		StatusCode: statusCode,
		Response: cm.StatusResponse{
			Status:  status,
			Message: cm.NewOptString(message),
		},
	}
}

func (e statusMessageError) Error() string {
	return fmt.Sprintf("%d - %s: %s", e.StatusCode, e.Response.Status, e.Response.Message.Value)
}
