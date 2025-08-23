package testing

import (
	"context"
	"net/http"
	"strconv"

	cm "github.com/cysp/terraform-provider-censusworkspace/internal/census-management-go"
)

func (h *Handler) CreateDataset(_ context.Context, req cm.CreateDatasetBody) (*cm.IdResponseStatusCode, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	// sourceID := req.SourceID

	h.datasetIDLast++
	modelID := h.datasetIDLast
	modelIDString := strconv.FormatInt(modelID, 10)

	model := NewDatasetFromCreateDatasetBody(modelID, req)

	h.Datasets[modelIDString] = &model

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

func (h *Handler) GetDataset(_ context.Context, params cm.GetDatasetParams) (*cm.DatasetResponseStatusCode, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	model, exists := h.Datasets[params.DatasetID]
	if !exists {
		return nil, errNotFound
	}

	return &cm.DatasetResponseStatusCode{
		StatusCode: http.StatusOK,
		Response: cm.DatasetResponse{
			Status: cm.ResponseStatusSuccess,
			Data:   *model,
		},
	}, nil
}

func (h *Handler) UpdateDataset(_ context.Context, req cm.UpdateDatasetBody, params cm.UpdateDatasetParams) (*cm.DatasetResponseStatusCode, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	model, exists := h.Datasets[params.DatasetID]
	if !exists {
		return nil, errNotFound
	}

	err := UpdateDatasetWithUpdateDatasetBody(model, req)
	if err != nil {
		return nil, err
	}

	return &cm.DatasetResponseStatusCode{
		StatusCode: http.StatusOK,
		Response: cm.DatasetResponse{
			Status: cm.ResponseStatusSuccess,
			Data:   *model,
		},
	}, nil
}

func (h *Handler) DeleteDataset(_ context.Context, params cm.DeleteDatasetParams) (*cm.StatusResponseStatusCode, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	_, exists := h.Datasets[params.DatasetID]
	if !exists {
		return nil, errNotFound
	}

	delete(h.Datasets, params.DatasetID)

	return &cm.StatusResponseStatusCode{
		StatusCode: http.StatusOK,
		Response: cm.StatusResponse{
			Status: cm.NewResponseStatusStatusResponseStatus(cm.ResponseStatusDeleted),
		},
	}, nil
}
