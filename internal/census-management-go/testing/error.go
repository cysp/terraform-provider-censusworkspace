package testing

import (
	"context"
	"errors"
	"net/http"

	cm "github.com/cysp/terraform-provider-censusworkspace/internal/census-management-go"
)

var errNotFound = errors.New("not found")

func (h *Handler) NewError(ctx context.Context, err error) *cm.StatusResponseStatusCode {
	if err == errNotFound {
		return &cm.StatusResponseStatusCode{
			StatusCode: http.StatusNotFound,
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
