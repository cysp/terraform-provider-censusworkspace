//nolint:dupl
package testing

import (
	"context"
	"fmt"
	"net/http"

	cm "github.com/cysp/terraform-provider-censusworkspace/internal/census-management-go"
)

func (h *Handler) CreateSource(ctx context.Context, req *cm.CreateSourceData) (*cm.SourceResponseStatusCode, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.sourceIDLast++
	id := h.sourceIDLast

	source := NewSourceFromCreateSourceData(id, *req)

	h.Sources[fmt.Sprintf("%d", id)] = &source

	return &cm.SourceResponseStatusCode{
		StatusCode: http.StatusCreated,
		Response: cm.SourceResponse{
			Status: cm.ResponseStatusSuccess,
			Data:   source,
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

func (h *Handler) UpdateSource(ctx context.Context, req *cm.UpdateSourceData, params cm.UpdateSourceParams) (*cm.SourceResponseStatusCode, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	source, exists := h.Sources[params.SourceID]
	if !exists {
		return nil, errNotFound
	}

	UpdateSourceWithUpdateSourceData(source, *req)

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
			Status:  cm.ResponseStatusSuccess,
			Message: cm.NewOptString("Source deleted"),
		},
	}, nil
}
