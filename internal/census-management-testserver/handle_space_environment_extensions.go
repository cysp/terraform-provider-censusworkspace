//nolint:dupl
package censusmanagementtestserver

import (
	"net/http"

	cm "github.com/cysp/terraform-provider-census/internal/census-management-go"
)

func (ts *CensusManagementTestServer) SetupSpaceEnvironmentExtensionHandlers() {
	ts.serveMux.Handle("/spaces/{spaceID}/environments/{environmentID}/extensions/{extensionID}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		spaceID := r.PathValue("spaceID")
		environmentID := r.PathValue("environmentID")
		extensionID := r.PathValue("extensionID")

		if spaceID == NonexistentID || environmentID == NonexistentID || extensionID == NonexistentID {
			_ = WriteCensusManagementErrorNotFoundResponse(w)

			return
		}

		ts.mu.Lock()
		defer ts.mu.Unlock()

		extension := ts.extensions.Get(spaceID, environmentID, extensionID)

		switch r.Method {
		case http.MethodGet:
			switch extension {
			case nil:
				_ = WriteCensusManagementErrorNotFoundResponse(w)
			default:
				_ = WriteCensusManagementResponse(w, http.StatusOK, extension)
			}

		case http.MethodPut:
			var extensionFields cm.ExtensionFields
			if err := ReadCensusManagementRequest(r, &extensionFields); err != nil {
				_ = WriteCensusManagementErrorBadRequestResponseWithError(w, err)

				return
			}

			switch extension {
			case nil:
				appInstallation := NewExtensionFromFields(spaceID, environmentID, extensionID, extensionFields)
				ts.extensions.Set(spaceID, environmentID, extensionID, &appInstallation)
				_ = WriteCensusManagementResponse(w, http.StatusOK, &appInstallation)
			default:
				UpdateExtensionFromFields(extension, extensionFields)
				_ = WriteCensusManagementResponse(w, http.StatusOK, extension)
			}

		case http.MethodDelete:
			switch extension {
			case nil:
				_ = WriteCensusManagementErrorNotFoundResponse(w)
			default:
				ts.extensions.Delete(spaceID, environmentID, extensionID)
				w.WriteHeader(http.StatusNoContent)
			}

		default:
			_ = WriteCensusManagementErrorNotFoundResponse(w)
		}
	}))
}
