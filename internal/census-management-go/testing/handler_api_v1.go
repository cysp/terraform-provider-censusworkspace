package testing

import (
	"context"

	cm "github.com/cysp/terraform-provider-censusworkspace/internal/census-management-go"
)

//nolint:ireturn,revive
func (h *Handler) GetApiV1(_ context.Context) (cm.GetApiV1Res, error) {
	return nil, errNotFound
}
