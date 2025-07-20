//nolint:dupl
package censusmanagementtestserver

import (
	"net/http"

	cm "github.com/cysp/terraform-provider-census/internal/census-management-go"
)

func (ts *CensusManagementTestServer) SetupSpaceEnvironmentAppInstallationHandlers() {
	//nolint:dupl
	ts.serveMux.Handle("/spaces/{spaceID}/environments/{environmentID}/app_installations/{appDefinitionID}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		spaceID := r.PathValue("spaceID")
		environmentID := r.PathValue("environmentID")
		appDefinitionID := r.PathValue("appDefinitionID")

		if spaceID == NonexistentID || environmentID == NonexistentID || appDefinitionID == NonexistentID {
			_ = WriteCensusManagementErrorNotFoundResponse(w)

			return
		}

		ts.mu.Lock()
		defer ts.mu.Unlock()

		appInstallation := ts.appInstallations.Get(spaceID, environmentID, appDefinitionID)

		switch r.Method {
		case http.MethodGet:
			switch appInstallation {
			case nil:
				_ = WriteCensusManagementErrorNotFoundResponse(w)
			default:
				_ = WriteCensusManagementResponse(w, http.StatusOK, appInstallation)
			}

		case http.MethodPut:
			var appInstallationFields cm.AppInstallationFields
			if err := ReadCensusManagementRequest(r, &appInstallationFields); err != nil {
				_ = WriteCensusManagementErrorBadRequestResponseWithError(w, err)

				return
			}

			switch appInstallation {
			case nil:
				appInstallation := NewAppInstallationFromFields(spaceID, environmentID, appDefinitionID, appInstallationFields)
				ts.appInstallations.Set(spaceID, environmentID, appDefinitionID, &appInstallation)
				_ = WriteCensusManagementResponse(w, http.StatusOK, &appInstallation)
			default:
				UpdateAppInstallationFromFields(appInstallation, appInstallationFields)
				_ = WriteCensusManagementResponse(w, http.StatusOK, appInstallation)
			}

		case http.MethodDelete:
			switch appInstallation {
			case nil:
				_ = WriteCensusManagementErrorNotFoundResponse(w)
			default:
				ts.appInstallations.Delete(spaceID, environmentID, appDefinitionID)
				w.WriteHeader(http.StatusNoContent)
			}

		default:
			_ = WriteCensusManagementErrorNotFoundResponse(w)
		}
	}))
}
