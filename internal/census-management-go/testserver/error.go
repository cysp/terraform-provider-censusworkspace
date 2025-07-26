package censusmanagementtestserver

import (
	"context"
	"errors"
	"net/http"

	cm "github.com/cysp/terraform-provider-censusworkspace/internal/census-management-go/server"
)

type NotFoundError struct {
}

func (e NotFoundError) Error() string {
	return "not found"
}

func (h *handler) NewError(ctx context.Context, err error) *cm.StatusResponseStatusCode {
	var notFoundError NotFoundError
	if errors.As(err, &notFoundError) {
		return &cm.StatusResponseStatusCode{
			StatusCode: http.StatusNotFound,
			Response: cm.StatusResponse{
				Status:  cm.ResponseStatusNotFound,
				Message: cm.NewOptString(notFoundError.Error()),
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
