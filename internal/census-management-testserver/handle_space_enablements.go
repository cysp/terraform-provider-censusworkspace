package censusmanagementtestserver

import (
	"net/http"

	cm "github.com/cysp/terraform-provider-census/internal/census-management-go"
)

func (ts *CensusManagementTestServer) setupSpaceEnablementsHandlers() {
	ts.serveMux.Handle("/spaces/{spaceID}/enablements", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		spaceID := r.PathValue("spaceID")

		if spaceID == NonexistentID {
			_ = WriteCensusManagementErrorNotFoundResponse(w)

			return
		}

		ts.mu.Lock()
		defer ts.mu.Unlock()

		enablements := ts.getOrCreateSpaceEnablements(spaceID)

		switch r.Method {
		case http.MethodGet:
			_ = WriteCensusManagementResponse(w, http.StatusOK, enablements)

		case http.MethodPut:
			var enablementRequestFields cm.SpaceEnablementFields
			if err := ReadCensusManagementRequest(r, &enablementRequestFields); err != nil {
				_ = WriteCensusManagementErrorBadRequestResponseWithError(w, err)

				return
			}

			UpdateSpaceEnablementFromRequestFields(enablements, enablementRequestFields)

			_ = WriteCensusManagementResponse(w, http.StatusOK, enablements)

		default:
			_ = WriteCensusManagementErrorNotFoundResponse(w)
		}
	}))
}

func (ts *CensusManagementTestServer) getOrCreateSpaceEnablements(spaceID string) *cm.SpaceEnablement {
	enablements, ok := ts.enablements[spaceID]
	if !ok {
		enablements = pointerTo(NewSpaceEnablement(spaceID))
		ts.enablements[spaceID] = enablements
	}

	return enablements
}

func (ts *CensusManagementTestServer) SetSpaceEnablements(spaceID string, enablementFields cm.SpaceEnablementFields) {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	spaceEnablement := NewSpaceEnablementFromRequestFields(spaceID, enablementFields)

	ts.enablements[spaceID] = &spaceEnablement
}
