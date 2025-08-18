package provider

import (
	"context"

	cm "github.com/cysp/terraform-provider-censusworkspace/internal/census-management-go"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func (m *SourceModel) ToCreateSourceData(_ context.Context) (cm.CreateSourceBody, diag.Diagnostics) {
	body := cm.CreateSourceBody{
		Connection: cm.CreateSourceBodyConnection{
			Type: m.Type.ValueString(),
		},
	}

	syncEngine := m.SyncEngine.ValueString()
	if syncEngine != "" {
		body.Connection.SyncEngine.SetTo(syncEngine)
	}

	label := m.Label.ValueString()
	if label != "" {
		body.Connection.Label.SetTo(label)
	} else {
		body.Connection.Label.SetToNull()
	}

	if credentials := m.Credentials.ValueStringPointer(); credentials != nil {
		body.Connection.Credentials = []byte(*credentials)
	}

	return body, nil
}

func (m *SourceModel) ToUpdateSourceData(_ context.Context) (cm.UpdateSourceBody, diag.Diagnostics) {
	body := cm.UpdateSourceBody{}

	label := m.Label.ValueString()
	if label != "" {
		body.Connection.Label.SetTo(label)
	} else {
		body.Connection.Label.SetToNull()
	}

	if credentials := m.Credentials.ValueStringPointer(); credentials != nil {
		body.Connection.Credentials = []byte(*credentials)
	}

	return body, nil
}
