//nolint:dupl
package testing

import (
	"context"
	"fmt"
	"net/http"

	cm "github.com/cysp/terraform-provider-censusworkspace/internal/census-management-go"
)

func (h *Handler) CreateSource(ctx context.Context, req *cm.CreateSourceBody) (*cm.IdResponseStatusCode, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.sourceIDLast++
	sourceID := h.sourceIDLast

	source := NewSourceFromCreateSourceBody(sourceID, *req)

	h.Sources[fmt.Sprintf("%d", sourceID)] = &source

	return &cm.IdResponseStatusCode{
		StatusCode: http.StatusCreated,
		Response: cm.IdResponse{
			Status: cm.ResponseStatusSuccess,
			Data: cm.IdResponseData{
				ID: sourceID,
			},
		},
	}, nil
}

func (h *Handler) GetSource(ctx context.Context, params cm.GetSourceParams) (*cm.SourceResponseStatusCode, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	source, exists := h.Sources[params.SourceID]
	if !exists {
		return nil, errNotFound
	}

	return &cm.SourceResponseStatusCode{
		StatusCode: http.StatusOK,
		Response: cm.SourceResponse{
			Status: cm.ResponseStatusSuccess,
			Data:   *source,
		},
	}, nil
}

func (h *Handler) UpdateSource(ctx context.Context, req *cm.UpdateSourceBody, params cm.UpdateSourceParams) (*cm.SourceResponseStatusCode, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	source, exists := h.Sources[params.SourceID]
	if !exists {
		return nil, errNotFound
	}

	UpdateSourceWithUpdateSourceBody(source, *req)

	return &cm.SourceResponseStatusCode{
		StatusCode: http.StatusOK,
		Response: cm.SourceResponse{
			Status: cm.ResponseStatusSuccess,
			Data:   *source,
		},
	}, nil
}

func (h *Handler) DeleteSource(ctx context.Context, params cm.DeleteSourceParams) (*cm.StatusResponseStatusCode, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	_, exists := h.Sources[params.SourceID]
	if !exists {
		return nil, errNotFound
	}

	delete(h.Sources, params.SourceID)

	return &cm.StatusResponseStatusCode{
		StatusCode: http.StatusOK,
		Response: cm.StatusResponse{
			Status: cm.NewResponseStatusStatusResponseStatus(cm.ResponseStatusDeleted),
		},
	}, nil
}
