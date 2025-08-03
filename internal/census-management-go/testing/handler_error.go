package testing

import (
	"context"
	"errors"
	"net/http"

	cm "github.com/cysp/terraform-provider-censusworkspace/internal/census-management-go"
)

func (h *Handler) NewError(ctx context.Context, err error) *cm.StatusResponseStatusCode {
	var sme statusMessageError
	if errors.As(err, &sme) {
		return &cm.StatusResponseStatusCode{
			StatusCode: sme.StatusCode,
			Response: cm.StatusResponse{
				Status:  cm.ResponseStatusNotFound,
				Message: cm.NewOptString(err.Error()),
			},
		}
	}

	return &cm.StatusResponseStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: cm.StatusResponse{
			Status:  cm.ResponseStatusError,
			Message: cm.NewOptString(err.Error()),
		},
	}
}
