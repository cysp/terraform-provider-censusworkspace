//nolint:dupl
package censusmanagementtestserver

import (
	"net/http"

	cm "github.com/cysp/terraform-provider-census/internal/census-management-go"
)

func (ts *CensusManagementTestServer) setupSpaceWebhookDefinitionHandlers() {
	ts.serveMux.Handle("/spaces/{spaceID}/webhook_definitions", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		spaceID := r.PathValue("spaceID")

		if spaceID == NonexistentID {
			_ = WriteCensusManagementErrorNotFoundResponse(w)

			return
		}

		ts.mu.Lock()
		defer ts.mu.Unlock()

		switch r.Method {
		case http.MethodPost:
			var webhookDefinitionFields cm.WebhookDefinitionFields
			if err := ReadCensusManagementRequest(r, &webhookDefinitionFields); err != nil {
				_ = WriteCensusManagementErrorBadRequestResponseWithError(w, err)

				return
			}

			webhookDefinition := NewWebhookDefinitionFromFields(spaceID, generateResourceID(), webhookDefinitionFields)

			ts.webhookDefinitions.Set(spaceID, webhookDefinition.Sys.ID, &webhookDefinition)

			_ = WriteCensusManagementResponse(w, http.StatusCreated, &webhookDefinition)

		default:
			_ = WriteCensusManagementErrorNotFoundResponse(w)
		}
	}))

	ts.serveMux.Handle("/spaces/{spaceID}/webhook_definitions/{webhookDefinitionID}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		spaceID := r.PathValue("spaceID")
		webhookDefinitionID := r.PathValue("webhookDefinitionID")

		if spaceID == NonexistentID || webhookDefinitionID == NonexistentID {
			_ = WriteCensusManagementErrorNotFoundResponse(w)

			return
		}

		ts.mu.Lock()
		defer ts.mu.Unlock()

		webhookDefinition := ts.webhookDefinitions.Get(spaceID, webhookDefinitionID)

		switch r.Method {
		case http.MethodGet:
			switch webhookDefinition {
			case nil:
				_ = WriteCensusManagementErrorNotFoundResponse(w)
			default:
				_ = WriteCensusManagementResponse(w, http.StatusOK, webhookDefinition)
			}

		case http.MethodPut:
			var webhookDefinitionFields cm.WebhookDefinitionFields
			if err := ReadCensusManagementRequest(r, &webhookDefinitionFields); err != nil {
				_ = WriteCensusManagementErrorBadRequestResponseWithError(w, err)

				return
			}

			switch webhookDefinition {
			case nil:
				webhookDefinition := NewWebhookDefinitionFromFields(spaceID, webhookDefinitionID, webhookDefinitionFields)
				ts.webhookDefinitions.Set(spaceID, webhookDefinition.Sys.ID, &webhookDefinition)
				_ = WriteCensusManagementResponse(w, http.StatusCreated, &webhookDefinition)
			default:
				UpdateWebhookDefinitionFromFields(webhookDefinition, webhookDefinitionFields)
				_ = WriteCensusManagementResponse(w, http.StatusOK, webhookDefinition)
			}

		case http.MethodDelete:
			switch webhookDefinition {
			case nil:
				_ = WriteCensusManagementErrorNotFoundResponse(w)
			default:
				ts.webhookDefinitions.Delete(spaceID, webhookDefinitionID)
				w.WriteHeader(http.StatusNoContent)
			}

		default:
			_ = WriteCensusManagementErrorNotFoundResponse(w)
		}
	}))
}

func (ts *CensusManagementTestServer) SetWebhookDefinition(spaceID, webhookDefinitionID string, webhookDefinitionFields cm.WebhookDefinitionFields) {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	webhookDefinition := NewWebhookDefinitionFromFields(spaceID, webhookDefinitionID, webhookDefinitionFields)

	ts.webhookDefinitions.Set(spaceID, webhookDefinition.Sys.ID, &webhookDefinition)
}
