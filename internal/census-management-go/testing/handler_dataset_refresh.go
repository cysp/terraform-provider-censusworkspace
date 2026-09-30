package testing

import (
	"context"
	"net/http"

	cm "github.com/cysp/terraform-provider-censusworkspace/internal/census-management-go"
)

// RefreshDatasetColumns implements censusmanagement.Handler.
func (h *Handler) RefreshDatasetColumns(_ context.Context, params cm.RefreshDatasetColumnsParams) (*cm.RefreshKeyResponseStatusCode, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	_, exists := h.Datasets[params.DatasetID]
	if !exists {
		return nil, errNotFound
	}

	h.datasetRefreshKeyLast++
	refreshKey := h.datasetRefreshKeyLast

	return &cm.RefreshKeyResponseStatusCode{
		StatusCode: http.StatusOK,
		Response: cm.RefreshKeyResponse{
			RefreshKey: refreshKey,
		},
	}, nil
}
