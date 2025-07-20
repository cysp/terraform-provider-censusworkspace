package censusmanagementtestserver

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func ReadCensusManagementRequest[T any](r *http.Request, v *T) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("failed to read body: %w", err)
	}
	defer r.Body.Close()

	err = json.Unmarshal(body, v)
	if err != nil {
		return fmt.Errorf("failed to unmarshal body: %w", err)
	}

	return nil
}

func ReadCensusManagementRequestWithValidation[T any](r *http.Request, v *T, validateFunc func(*T) error) error {
	if err := ReadCensusManagementRequest(r, v); err != nil {
		return err
	}

	if err := validateFunc(v); err != nil {
		return err
	}

	return nil
}
