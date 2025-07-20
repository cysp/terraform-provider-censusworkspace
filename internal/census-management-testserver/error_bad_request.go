package censusmanagementtestserver

import (
	"net/http"
)

func WriteCensusManagementErrorBadRequestResponse(w http.ResponseWriter) error {
	return WriteCensusManagementErrorResponse(w, http.StatusBadRequest, "BadRequest", pointerTo("The request was malformed or contained invalid parameters."), nil)
}

func WriteCensusManagementErrorBadRequestResponseWithDetails(w http.ResponseWriter, details string) error {
	return WriteCensusManagementErrorResponse(w, http.StatusBadRequest, "BadRequest", pointerTo("The request was malformed or contained invalid parameters."), []byte(details))
}

func WriteCensusManagementErrorBadRequestResponseWithError(w http.ResponseWriter, err error) error {
	return WriteCensusManagementErrorBadRequestResponseWithDetails(w, err.Error())
}
