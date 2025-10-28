//nolint:dupl
package testing

import (
	"context"
	"net/http"
	"strconv"

	cm "github.com/cysp/terraform-provider-censusworkspace/internal/census-management-go"
)

func (h *Handler) CreateSync(_ context.Context, req *cm.CreateSyncBody) (*cm.SyncIdResponseStatusCode, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.syncIDLast++
	syncID := h.syncIDLast

	sync := NewSyncFromCreateSyncBody(syncID, *req)

	h.Syncs[strconv.FormatInt(syncID, 10)] = &sync

	return &cm.SyncIdResponseStatusCode{
		StatusCode: http.StatusCreated,
		Response: cm.SyncIdResponse{
			Status: cm.ResponseStatusSuccess,
			Data: cm.SyncIdResponseData{
				SyncID: syncID,
			},
		},
	}, nil
}

func (h *Handler) GetSync(_ context.Context, params cm.GetSyncParams) (*cm.SyncResponseStatusCode, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	sync, exists := h.Syncs[params.SyncID]
	if !exists {
		return nil, errNotFound
	}

	return &cm.SyncResponseStatusCode{
		StatusCode: http.StatusOK,
		Response: cm.SyncResponse{
			Status: cm.ResponseStatusSuccess,
			Data:   *sync,
		},
	}, nil
}

func (h *Handler) UpdateSync(_ context.Context, req *cm.UpdateSyncBody, params cm.UpdateSyncParams) (*cm.SyncResponseStatusCode, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	sync, exists := h.Syncs[params.SyncID]
	if !exists {
		return nil, errNotFound
	}

	UpdateSyncWithUpdateSyncBody(sync, *req)

	return &cm.SyncResponseStatusCode{
		StatusCode: http.StatusOK,
		Response: cm.SyncResponse{
			Status: cm.ResponseStatusSuccess,
			Data:   *sync,
		},
	}, nil
}

func (h *Handler) DeleteSync(_ context.Context, params cm.DeleteSyncParams) (*cm.StatusResponseStatusCode, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	_, exists := h.Syncs[params.SyncID]
	if !exists {
		return nil, errNotFound
	}

	delete(h.Syncs, params.SyncID)

	return &cm.StatusResponseStatusCode{
		StatusCode: http.StatusOK,
		Response: cm.StatusResponse{
			Status: cm.NewResponseStatusStatusResponseStatus(cm.ResponseStatusDeleted),
		},
	}, nil
}
