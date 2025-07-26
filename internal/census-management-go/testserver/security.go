package testserver

import (
	"context"

	cm "github.com/cysp/terraform-provider-censusworkspace/internal/census-management-go"
)

func (h *handler) HandleWorkspaceApiKey(ctx context.Context, operationName cm.OperationName, t cm.WorkspaceApiKey) (context.Context, error) {
	return ctx, nil
}
