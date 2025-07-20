//nolint:dupl
package censusmanagementtestserver

import (
	"net/http"

	cm "github.com/cysp/terraform-provider-census/internal/census-management-go"
)

func (ts *CensusManagementTestServer) setupSpaceRoleHandlers() {
	ts.serveMux.Handle("/spaces/{spaceID}/roles", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		spaceID := r.PathValue("spaceID")

		if spaceID == NonexistentID {
			_ = WriteCensusManagementErrorNotFoundResponse(w)

			return
		}

		ts.mu.Lock()
		defer ts.mu.Unlock()

		switch r.Method {
		case http.MethodPost:
			var roleFields cm.RoleFields
			if err := ReadCensusManagementRequest(r, &roleFields); err != nil {
				_ = WriteCensusManagementErrorBadRequestResponseWithError(w, err)

				return
			}

			role := NewRoleFromFields(spaceID, generateResourceID(), roleFields)

			ts.roles.Set(spaceID, role.Sys.ID, &role)

			_ = WriteCensusManagementResponse(w, http.StatusCreated, &role)

		default:
			_ = WriteCensusManagementErrorNotFoundResponse(w)
		}
	}))

	ts.serveMux.Handle("/spaces/{spaceID}/roles/{roleID}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		spaceID := r.PathValue("spaceID")
		roleID := r.PathValue("roleID")

		if spaceID == NonexistentID || roleID == NonexistentID {
			_ = WriteCensusManagementErrorNotFoundResponse(w)

			return
		}

		ts.mu.Lock()
		defer ts.mu.Unlock()

		role := ts.roles.Get(spaceID, roleID)

		switch r.Method {
		case http.MethodGet:
			switch role {
			case nil:
				_ = WriteCensusManagementErrorNotFoundResponse(w)
			default:
				_ = WriteCensusManagementResponse(w, http.StatusOK, role)
			}

		case http.MethodPut:
			var roleFields cm.RoleFields
			if err := ReadCensusManagementRequest(r, &roleFields); err != nil {
				_ = WriteCensusManagementErrorBadRequestResponseWithError(w, err)

				return
			}

			switch role {
			case nil:
				role := NewRoleFromFields(spaceID, roleID, roleFields)
				ts.roles.Set(spaceID, role.Sys.ID, &role)
				_ = WriteCensusManagementResponse(w, http.StatusCreated, &role)
			default:
				UpdateRoleFromFields(role, roleFields)
				_ = WriteCensusManagementResponse(w, http.StatusOK, role)
			}

		case http.MethodDelete:
			switch role {
			case nil:
				_ = WriteCensusManagementErrorNotFoundResponse(w)
			default:
				ts.roles.Delete(spaceID, roleID)
				w.WriteHeader(http.StatusNoContent)
			}

		default:
			_ = WriteCensusManagementErrorNotFoundResponse(w)
		}
	}))
}

func (ts *CensusManagementTestServer) SetRole(spaceID, roleID string, roleFields cm.RoleFields) {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	role := NewRoleFromFields(spaceID, roleID, roleFields)

	ts.roles.Set(spaceID, role.Sys.ID, &role)
}
