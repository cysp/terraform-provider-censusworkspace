package testing

import (
	"context"
	"errors"
	"net/http"

	cm "github.com/cysp/terraform-provider-censusworkspace/internal/census-management-go"
)

func (h *Handler) NewError(_ context.Context, err error) *cm.StatusResponseStatusCode {
	var sce statusCodeError
	if errors.As(err, &sce) {
		return &cm.StatusResponseStatusCode{
			StatusCode: sce.StatusCode,
			Response: cm.StatusResponse{
				Status: cm.NewIntStatusResponseStatus(sce.StatusCode),
			},
		}
	}

	return &cm.StatusResponseStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: cm.StatusResponse{
			Status:  cm.NewResponseStatusStatusResponseStatus(cm.ResponseStatusError),
			Message: cm.NewOptString(err.Error()),
		},
	}
}
