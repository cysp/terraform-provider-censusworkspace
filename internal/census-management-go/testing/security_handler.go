package testing

import (
	"context"
	"sync"

	cm "github.com/cysp/terraform-provider-censusworkspace/internal/census-management-go"
)

type SecurityHandler struct {
	mu sync.Mutex
}

var _ cm.SecurityHandler = (*SecurityHandler)(nil)

func (h *SecurityHandler) HandleWorkspaceApiKey(ctx context.Context, operationName cm.OperationName, t cm.WorkspaceApiKey) (context.Context, error) {
	return ctx, nil
}
