//nolint:dupl
package testing

import (
	"context"
	"net/http"
	"strconv"

	cm "github.com/cysp/terraform-provider-censusworkspace/internal/census-management-go"
)

func (h *Handler) CreateDestination(_ context.Context, req *cm.CreateDestinationBody) (*cm.IdResponseStatusCode, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.destinationIDLast++
	destinationID := h.destinationIDLast

	destination := NewDestinationFromCreateDestinationBody(destinationID, *req)

	h.Destinations[strconv.FormatInt(destinationID, 10)] = &destination

	return &cm.IdResponseStatusCode{
		StatusCode: http.StatusCreated,
		Response: cm.IdResponse{
			Status: cm.ResponseStatusSuccess,
			Data: cm.IdResponseData{
				ID: destinationID,
			},
		},
	}, nil
}

func (h *Handler) GetDestination(_ context.Context, params cm.GetDestinationParams) (*cm.DestinationResponseStatusCode, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	destination, exists := h.Destinations[params.DestinationID]
	if !exists {
		return nil, errNotFound
	}

	return &cm.DestinationResponseStatusCode{
		StatusCode: http.StatusOK,
		Response: cm.DestinationResponse{
			Status: cm.ResponseStatusSuccess,
			Data:   *destination,
		},
	}, nil
}

func (h *Handler) UpdateDestination(_ context.Context, req *cm.UpdateDestinationBody, params cm.UpdateDestinationParams) (*cm.DestinationResponseStatusCode, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	destination, exists := h.Destinations[params.DestinationID]
	if !exists {
		return nil, errNotFound
	}

	UpdateDestinationWithUpdateDestinationBody(destination, *req)

	return &cm.DestinationResponseStatusCode{
		StatusCode: http.StatusOK,
		Response: cm.DestinationResponse{
			Status: cm.ResponseStatusSuccess,
			Data:   *destination,
		},
	}, nil
}

func (h *Handler) DeleteDestination(_ context.Context, params cm.DeleteDestinationParams) (*cm.StatusResponseStatusCode, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	_, exists := h.Destinations[params.DestinationID]
	if !exists {
		return nil, errNotFound
	}

	delete(h.Destinations, params.DestinationID)

	return &cm.StatusResponseStatusCode{
		StatusCode: http.StatusOK,
		Response: cm.StatusResponse{
			Status: cm.NewResponseStatusStatusResponseStatus(cm.ResponseStatusDeleted),
		},
	}, nil
}
