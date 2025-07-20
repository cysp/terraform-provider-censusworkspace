package censusmanagementtestserver

import (
	"net/http"

	cm "github.com/cysp/terraform-provider-census/internal/census-management-go"
)

//nolint:cyclop,gocognit
func (ts *CensusManagementTestServer) setupSpaceEnvironmentContentTypeHandlers() {
	//nolint:dupl
	ts.serveMux.Handle("/spaces/{spaceID}/environments/{environmentID}/content_types/{contentTypeID}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		spaceID := r.PathValue("spaceID")
		environmentID := r.PathValue("environmentID")
		contentTypeID := r.PathValue("contentTypeID")

		if spaceID == NonexistentID || environmentID == NonexistentID || contentTypeID == NonexistentID {
			_ = WriteCensusManagementErrorNotFoundResponse(w)

			return
		}

		ts.mu.Lock()
		defer ts.mu.Unlock()

		contentType := ts.contentTypes.Get(spaceID, environmentID, contentTypeID)

		switch r.Method {
		case http.MethodGet:
			switch contentType {
			case nil:
				_ = WriteCensusManagementErrorNotFoundResponse(w)
			default:
				_ = WriteCensusManagementResponse(w, http.StatusOK, contentType)
			}
		case http.MethodPut:
			var contentTypeRequestFields cm.ContentTypeRequestFields
			if err := ReadCensusManagementRequest(r, &contentTypeRequestFields); err != nil {
				_ = WriteCensusManagementErrorBadRequestResponseWithError(w, err)

				return
			}

			switch contentType {
			case nil:
				contentType := NewContentTypeFromRequestFields(spaceID, environmentID, contentTypeID, contentTypeRequestFields)
				ts.contentTypes.Set(spaceID, environmentID, contentTypeID, &contentType)
				_ = WriteCensusManagementResponse(w, http.StatusCreated, &contentType)
			default:
				UpdateContentTypeFromRequestFields(contentType, contentTypeRequestFields)
				_ = WriteCensusManagementResponse(w, http.StatusOK, contentType)
			}

		case http.MethodDelete:
			switch contentType {
			case nil:
				_ = WriteCensusManagementErrorNotFoundResponse(w)
			default:
				ts.contentTypes.Delete(spaceID, environmentID, contentTypeID)
				w.WriteHeader(http.StatusNoContent)
			}

		default:
			_ = WriteCensusManagementErrorNotFoundResponse(w)
		}
	}))

	ts.serveMux.Handle("/spaces/{spaceID}/environments/{environmentID}/content_types/{contentTypeID}/published", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		spaceID := r.PathValue("spaceID")
		environmentID := r.PathValue("environmentID")
		contentTypeID := r.PathValue("contentTypeID")

		if spaceID == NonexistentID || environmentID == NonexistentID || contentTypeID == NonexistentID {
			_ = WriteCensusManagementErrorNotFoundResponse(w)

			return
		}

		ts.mu.Lock()
		defer ts.mu.Unlock()

		contentType := ts.contentTypes.Get(spaceID, environmentID, contentTypeID)

		switch r.Method {
		case http.MethodPut:
			switch contentType {
			case nil:
				_ = WriteCensusManagementErrorNotFoundResponse(w)
			default:
				publishContentType(contentType)
				_ = WriteCensusManagementResponse(w, http.StatusOK, contentType)
			}

		case http.MethodDelete:
			switch contentType {
			case nil:
				_ = WriteCensusManagementErrorNotFoundResponse(w)
			default:
				contentType.Sys.PublishedVersion.Reset()
				w.WriteHeader(http.StatusNoContent)
			}

		default:
			_ = WriteCensusManagementErrorNotFoundResponse(w)
		}
	}))

	ts.serveMux.Handle("/spaces/{spaceID}/environments/{environmentID}/content_types/{contentTypeID}/editor_interface", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		spaceID := r.PathValue("spaceID")
		environmentID := r.PathValue("environmentID")
		contentTypeID := r.PathValue("contentTypeID")

		if spaceID == NonexistentID || environmentID == NonexistentID || contentTypeID == NonexistentID {
			_ = WriteCensusManagementErrorNotFoundResponse(w)

			return
		}

		ts.mu.Lock()
		defer ts.mu.Unlock()

		contentType := ts.contentTypes.Get(spaceID, environmentID, contentTypeID)
		editorInterface := ts.editorInterfaces.Get(spaceID, environmentID, contentTypeID)

		switch r.Method {
		case http.MethodGet:
			switch editorInterface {
			case nil:
				_ = WriteCensusManagementErrorNotFoundResponse(w)
			default:
				_ = WriteCensusManagementResponse(w, http.StatusOK, editorInterface)
			}

		case http.MethodPut:
			if contentType == nil {
				_ = WriteCensusManagementErrorNotFoundResponse(w)

				return
			}

			editorInterfaceFields := cm.EditorInterfaceFields{}
			if err := ReadCensusManagementRequest(r, &editorInterfaceFields); err != nil {
				_ = WriteCensusManagementErrorBadRequestResponseWithError(w, err)

				return
			}

			switch editorInterface {
			case nil:
				editorInterface := NewEditorInterfaceFromFields(spaceID, environmentID, contentTypeID, editorInterfaceFields)
				ts.editorInterfaces.Set(spaceID, environmentID, contentTypeID, &editorInterface)
				_ = WriteCensusManagementResponse(w, http.StatusOK, &editorInterface)
			default:
				UpdateEditorInterfaceFromFields(editorInterface, editorInterfaceFields)
				_ = WriteCensusManagementResponse(w, http.StatusOK, editorInterface)
			}

		default:
			_ = WriteCensusManagementErrorNotFoundResponse(w)
		}
	}))
}

func (ts *CensusManagementTestServer) SetContentType(spaceID, environmentID, contentTypeID string, contentTypeFields cm.ContentTypeRequestFields) {
	contentType := NewContentTypeFromRequestFields(spaceID, environmentID, contentTypeID, contentTypeFields)
	ts.contentTypes.Set(spaceID, environmentID, contentType.Sys.ID, &contentType)
}

func (ts *CensusManagementTestServer) SetEditorInterface(spaceID, environmentID, contentTypeID string, editorInterfaceFields cm.EditorInterfaceFields) {
	editorInterface := NewEditorInterfaceFromFields(spaceID, environmentID, contentTypeID, editorInterfaceFields)
	ts.editorInterfaces.Set(spaceID, environmentID, contentTypeID, &editorInterface)
}
