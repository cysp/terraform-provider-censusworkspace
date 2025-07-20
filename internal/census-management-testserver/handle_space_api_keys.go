package censusmanagementtestserver

import (
	"net/http"

	cm "github.com/cysp/terraform-provider-census/internal/census-management-go"
)

func (ts *CensusManagementTestServer) setupSpaceAPIKeyHandlers() {
	ts.serveMux.Handle("/spaces/{spaceID}/api_keys", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		spaceID := r.PathValue("spaceID")

		if spaceID == NonexistentID {
			_ = WriteCensusManagementErrorNotFoundResponse(w)

			return
		}

		ts.mu.Lock()
		defer ts.mu.Unlock()

		switch r.Method {
		case http.MethodPost:
			var apiKeyRequestFields cm.ApiKeyRequestFields
			if err := ReadCensusManagementRequest(r, &apiKeyRequestFields); err != nil {
				_ = WriteCensusManagementErrorBadRequestResponseWithError(w, err)

				return
			}

			apiKeyID := generateResourceID()
			apiKey := NewAPIKeyFromRequestFields(spaceID, apiKeyID, apiKeyRequestFields)
			apiKey.ApiKey = generateResourceID()

			previewAPIKeyID := generateResourceID()
			previewAPIKey := cm.PreviewApiKey{
				Sys:    NewPreviewAPIKeySys(spaceID, previewAPIKeyID),
				ApiKey: generateResourceID(),
			}

			apiKey.PreviewAPIKey.SetTo(cm.ApiKeyPreviewAPIKey{
				Sys: cm.ApiKeyPreviewAPIKeySys{
					Type:     cm.ApiKeyPreviewAPIKeySysTypeLink,
					LinkType: cm.ApiKeyPreviewAPIKeySysLinkTypePreviewApiKey,
					ID:       previewAPIKeyID,
				},
			})

			ts.apiKeys.Set(spaceID, apiKeyID, &apiKey)

			ts.previewAPIKeys.Set(spaceID, previewAPIKeyID, &previewAPIKey)

			_ = WriteCensusManagementResponse(w, http.StatusCreated, &apiKey)

		default:
			_ = WriteCensusManagementErrorNotFoundResponse(w)
		}
	}))

	ts.serveMux.Handle("/spaces/{spaceID}/api_keys/{apiKeyID}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		spaceID := r.PathValue("spaceID")
		apiKeyID := r.PathValue("apiKeyID")

		if spaceID == NonexistentID || apiKeyID == NonexistentID {
			_ = WriteCensusManagementErrorNotFoundResponse(w)

			return
		}

		ts.mu.Lock()
		defer ts.mu.Unlock()

		apiKey := ts.apiKeys.Get(spaceID, apiKeyID)

		switch r.Method {
		case http.MethodGet:
			switch apiKey {
			case nil:
				_ = WriteCensusManagementErrorNotFoundResponse(w)
			default:
				_ = WriteCensusManagementResponse(w, http.StatusOK, apiKey)
			}

		case http.MethodPut:
			var apiKeyRequestFields cm.ApiKeyRequestFields
			if err := ReadCensusManagementRequest(r, &apiKeyRequestFields); err != nil {
				_ = WriteCensusManagementErrorBadRequestResponseWithError(w, err)

				return
			}

			switch apiKey {
			case nil:
				apiKey := NewAPIKeyFromRequestFields(spaceID, apiKeyID, apiKeyRequestFields)
				ts.apiKeys.Set(spaceID, apiKeyID, &apiKey)
				_ = WriteCensusManagementResponse(w, http.StatusCreated, &apiKey)
			default:
				UpdateAPIKeyFromRequestFields(apiKey, apiKeyRequestFields)
				_ = WriteCensusManagementResponse(w, http.StatusOK, apiKey)
			}

		case http.MethodDelete:
			switch apiKey {
			case nil:
				_ = WriteCensusManagementErrorNotFoundResponse(w)
			default:
				ts.apiKeys.Delete(spaceID, apiKeyID)
				w.WriteHeader(http.StatusNoContent)
			}

		default:
			_ = WriteCensusManagementErrorNotFoundResponse(w)
		}
	}))
}
