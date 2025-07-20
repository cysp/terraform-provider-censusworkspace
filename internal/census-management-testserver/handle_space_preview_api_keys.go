package censusmanagementtestserver

import (
	"net/http"
)

func (ts *CensusManagementTestServer) setupSpacePreviewAPIKeyHandlers() {
	ts.serveMux.Handle("/spaces/{spaceID}/preview_api_keys/{previewAPIKeyID}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		spaceID := r.PathValue("spaceID")
		previewAPIKeyID := r.PathValue("previewAPIKeyID")

		if spaceID == NonexistentID || previewAPIKeyID == NonexistentID {
			_ = WriteCensusManagementErrorNotFoundResponse(w)

			return
		}

		ts.mu.Lock()
		defer ts.mu.Unlock()

		previewAPIKey := ts.previewAPIKeys.Get(spaceID, previewAPIKeyID)

		switch r.Method {
		case http.MethodGet:
			switch previewAPIKey {
			case nil:
				_ = WriteCensusManagementErrorNotFoundResponse(w)
			default:
				_ = WriteCensusManagementResponse(w, http.StatusOK, previewAPIKey)
			}

		default:
			_ = WriteCensusManagementErrorNotFoundResponse(w)
		}
	}))
}
