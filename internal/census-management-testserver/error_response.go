package censusmanagementtestserver

import (
	"net/http"

	cm "github.com/cysp/terraform-provider-census/internal/census-management-go"
)

func WriteCensusManagementErrorResponse(w http.ResponseWriter, statusCode int, id string, message *string, details []byte) error {
	return WriteCensusManagementResponse(w, statusCode, &cm.Error{
		Sys: cm.ErrorSys{
			Type: cm.ErrorSysTypeError,
			ID:   id,
		},
		Message: cm.NewOptPointerString(message),
		Details: details,
	})
}
