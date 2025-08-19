package censusmanagement

import (
	"context"
)

type WorkspaceAPIKeySecuritySource struct {
	token string
}

var _ SecuritySource = (*WorkspaceAPIKeySecuritySource)(nil)

func NewWorkspaceAPIKeySecuritySource(token string) WorkspaceAPIKeySecuritySource {
	return WorkspaceAPIKeySecuritySource{
		token: token,
	}
}

//nolint:revive
func (c WorkspaceAPIKeySecuritySource) WorkspaceApiKey(_ context.Context, _ string, _ *Client) (WorkspaceApiKey, error) {
	return WorkspaceApiKey{Token: c.token}, nil
}
