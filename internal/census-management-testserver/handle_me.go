package censusmanagementtestserver

import (
	"net/http"

	cm "github.com/cysp/terraform-provider-census/internal/census-management-go"
)

func (ts *CensusManagementTestServer) setupUserHandler() {
	ts.serveMux.Handle("/users/me", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ts.mu.Lock()
		defer ts.mu.Unlock()

		switch r.Method {
		case http.MethodGet:
			switch ts.me {
			case nil:
				_ = WriteCensusManagementErrorNotFoundResponse(w)
			default:
				_ = WriteCensusManagementResponse(w, http.StatusOK, ts.me)
			}
		default:
			_ = WriteCensusManagementErrorNotFoundResponse(w)
		}
	}))
}

func (ts *CensusManagementTestServer) SetMe(me *cm.User) {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	ts.me = me
}
