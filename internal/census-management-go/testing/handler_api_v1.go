package testing

import (
	"context"

	cm "github.com/cysp/terraform-provider-censusworkspace/internal/census-management-go"
)

func (h *Handler) GetApiV1(ctx context.Context) (cm.GetApiV1Res, error) {
	return nil, errNotFound
}
