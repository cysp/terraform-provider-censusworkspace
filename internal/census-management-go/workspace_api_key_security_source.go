package censusmanagement

import (
	"context"
)

type WorkspaceApiKeySecuritySource struct {
	token string
}

var _ SecuritySource = (*WorkspaceApiKeySecuritySource)(nil)

func NewWorkspaceApiKeySecuritySource(token string) WorkspaceApiKeySecuritySource {
	return WorkspaceApiKeySecuritySource{
		token: token,
	}
}

func (c WorkspaceApiKeySecuritySource) WorkspaceApiKey(_ context.Context, _ string, _ *Client) (WorkspaceApiKey, error) {
	return WorkspaceApiKey{Token: c.token}, nil
}
