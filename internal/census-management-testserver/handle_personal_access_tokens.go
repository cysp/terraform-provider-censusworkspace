package censusmanagementtestserver

import (
	"net/http"
	"time"

	cm "github.com/cysp/terraform-provider-census/internal/census-management-go"
)

func (ts *CensusManagementTestServer) setupPersonalAccessTokenHandlers() {
	ts.serveMux.Handle("/users/me/api_keys", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ts.mu.Lock()
		defer ts.mu.Unlock()

		switch r.Method {
		case http.MethodPost:
			var personalAccessTokenRequestFields cm.PersonalAccessTokenRequestFields
			if err := ReadCensusManagementRequestWithValidation(r, &personalAccessTokenRequestFields, validatePersonalAccessTokenRequestFields); err != nil {
				_ = WriteCensusManagementErrorBadRequestResponseWithError(w, err)

				return
			}

			personalAccessTokenID := generateResourceID()
			personalAccessToken := NewPersonalAccessTokenFromRequestFields(personalAccessTokenID, personalAccessTokenRequestFields)

			ts.personalAccessTokens[personalAccessToken.Sys.ID] = &personalAccessToken

			personalAccessTokenWithToken := personalAccessToken
			personalAccessTokenWithToken.Token.SetTo(generateResourceID())

			_ = WriteCensusManagementResponse(w, http.StatusCreated, &personalAccessTokenWithToken)

		default:
			_ = WriteCensusManagementErrorNotFoundResponse(w)
		}
	}))

	ts.serveMux.Handle("/users/me/api_keys/{id}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id") //nolint:varnamelen

		if id == NonexistentID {
			_ = WriteCensusManagementErrorNotFoundResponse(w)

			return
		}

		ts.mu.Lock()
		defer ts.mu.Unlock()

		personalAccessToken := ts.personalAccessTokens[id]

		switch r.Method {
		case http.MethodGet:
			switch personalAccessToken {
			case nil:
				_ = WriteCensusManagementErrorNotFoundResponse(w)

			default:
				_ = WriteCensusManagementResponse(w, http.StatusOK, personalAccessToken)
			}

		default:
			_ = WriteCensusManagementErrorNotFoundResponse(w)
		}
	}))

	ts.serveMux.Handle("/users/me/api_keys/{id}/revoked", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id") //nolint:varnamelen

		if id == NonexistentID {
			_ = WriteCensusManagementErrorNotFoundResponse(w)

			return
		}

		ts.mu.Lock()
		defer ts.mu.Unlock()

		personalAccessToken := ts.personalAccessTokens[id]

		switch r.Method {
		case http.MethodPut:
			switch personalAccessToken {
			case nil:
				_ = WriteCensusManagementErrorNotFoundResponse(w)
			default:
				personalAccessToken.RevokedAt.SetTo(time.Now())
				_ = WriteCensusManagementResponse(w, http.StatusOK, personalAccessToken)
			}

		default:
			_ = WriteCensusManagementErrorNotFoundResponse(w)
		}
	}))
}
