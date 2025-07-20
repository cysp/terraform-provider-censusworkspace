package censusmanagementtestserver

import (
	"net/http"
)

func WriteCensusManagementErrorNotFoundResponse(w http.ResponseWriter) error {
	return WriteCensusManagementErrorResponse(w, http.StatusNotFound, "NotFound", pointerTo("The resource could not be found."), nil)
}
