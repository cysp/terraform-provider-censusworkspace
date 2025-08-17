package testing

import (
	"context"
	"fmt"
	"net/http"

	cm "github.com/cysp/terraform-provider-censusworkspace/internal/census-management-go"
)

func (h *Handler) CreateSourceModel(ctx context.Context, req *cm.CreateSourceModelBody, params cm.CreateSourceModelParams) (*cm.IdResponseStatusCode, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	sourceID := params.SourceID

	h.modelIDLast++
	modelID := h.modelIDLast
	modelIDString := fmt.Sprintf("%d", modelID)

	model := NewSourceModelFromCreateSourceModelBody(modelID, *req)

	h.Models.Set(sourceID, modelIDString, &model)

	return &cm.IdResponseStatusCode{
		StatusCode: http.StatusCreated,
		Response: cm.IdResponse{
			Status: cm.ResponseStatusSuccess,
			Data: cm.IdResponseData{
				ID: modelID,
			},
		},
	}, nil
}

func (h *Handler) GetSourceModel(ctx context.Context, params cm.GetSourceModelParams) (*cm.SourceModelResponseStatusCode, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	model, exists := h.Models.Get(params.SourceID, params.ModelID)
	if !exists {
		return nil, errNotFound
	}

	return &cm.SourceModelResponseStatusCode{
		StatusCode: http.StatusOK,
		Response: cm.SourceModelResponse{
			Status: cm.ResponseStatusSuccess,
			Data:   *model,
		},
	}, nil
}

func (h *Handler) UpdateSourceModel(ctx context.Context, req *cm.UpdateSourceModelBody, params cm.UpdateSourceModelParams) (*cm.SourceModelResponseStatusCode, error) {
	panic("unimplemented")
}

func (h *Handler) DeleteSourceModel(ctx context.Context, params cm.DeleteSourceModelParams) (*cm.StatusResponseStatusCode, error) {
	panic("unimplemented")
}
