package censusmanagementtestserver

import (
	"context"

	cm "github.com/cysp/terraform-provider-censusworkspace/internal/census-management-go/server"
)

func (h *handler) HandleWorkspaceApiKey(ctx context.Context, operationName cm.OperationName, t cm.WorkspaceApiKey) (context.Context, error) {
	return ctx, nil
}
