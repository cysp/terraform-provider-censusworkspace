//nolint:dupl
package censusmanagementtestserver

import (
	"context"
	"errors"
	"net/http"

	cm "github.com/cysp/terraform-provider-censusworkspace/internal/census-management-go/server"
)

func (h *handler) CreateSource(ctx context.Context, req *cm.CreateSourceData) (*cm.SourceResponseStatusCode, error) {
	h.ts.mu.Lock()
	defer h.ts.mu.Unlock()

	h.ts.sourceIDLast++
	id := h.ts.sourceIDLast

	source := NewSourceFromCreateSourceData(id, *req)

	h.ts.Sources[id] = &source

	return &cm.SourceResponseStatusCode{
		StatusCode: http.StatusCreated,
		Response: cm.SourceResponse{
			Status: cm.ResponseStatusSuccess,
			Data:   source,
		},
	}, nil
}

func (h *handler) GetSource(ctx context.Context, params cm.GetSourceParams) (*cm.SourceResponseStatusCode, error) {
	h.ts.mu.Lock()
	defer h.ts.mu.Unlock()

	source, exists := h.ts.Sources[params.SourceID]
	if !exists {
		return nil, errors.New("not found")
	}

	return &cm.SourceResponseStatusCode{
		StatusCode: http.StatusOK,
		Response: cm.SourceResponse{
			Status: cm.ResponseStatusSuccess,
			Data:   *source,
		},
	}, nil
}

func (h *handler) UpdateSource(ctx context.Context, req *cm.UpdateSourceData, params cm.UpdateSourceParams) (*cm.SourceResponseStatusCode, error) {
	h.ts.mu.Lock()
	defer h.ts.mu.Unlock()

	source, exists := h.ts.Sources[params.SourceID]
	if !exists {
		return nil, errors.New("not found")
	}

	UpdateSourceFromUpdateSourceData(source, *req)

	return &cm.SourceResponseStatusCode{
		StatusCode: http.StatusOK,
		Response: cm.SourceResponse{
			Status: cm.ResponseStatusSuccess,
			Data:   *source,
		},
	}, nil
}

func (h *handler) DeleteSource(ctx context.Context, params cm.DeleteSourceParams) (*cm.StatusResponseStatusCode, error) {
	h.ts.mu.Lock()
	defer h.ts.mu.Unlock()

	_, exists := h.ts.Sources[params.SourceID]
	if !exists {
		return nil, errors.New("not found")
	}

	delete(h.ts.Sources, params.SourceID)

	return &cm.StatusResponseStatusCode{
		StatusCode: http.StatusNoContent,
		Response: cm.StatusResponse{
			Status:  cm.ResponseStatusSuccess,
			Message: cm.NewOptString("Source deleted"),
		},
	}, nil
}
