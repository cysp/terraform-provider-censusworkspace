package provider

import (
	"context"

	cm "github.com/cysp/terraform-provider-censusworkspace/internal/census-management-go"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func (m *SourceModel) ToCreateSourceData(_ context.Context) (cm.CreateSourceBody, diag.Diagnostics) {
	fields := cm.CreateSourceBody{
		Connection: cm.CreateSourceBodyConnection{
			Type:  m.Type.ValueString(),
			Label: cm.NewNilPointerString(m.Label.ValueStringPointer()),
		},
	}

	if credentials := m.Credentials.ValueStringPointer(); credentials != nil {
		fields.Connection.Credentials = []byte(*credentials)
	}

	return fields, nil
}

func (m *SourceModel) ToUpdateSourceData(_ context.Context) (cm.UpdateSourceBody, diag.Diagnostics) {
	fields := cm.UpdateSourceBody{
		Connection: cm.UpdateSourceBodyConnection{
			Label: cm.NewNilPointerString(m.Label.ValueStringPointer()),
		},
	}

	if credentials := m.Credentials.ValueStringPointer(); credentials != nil {
		fields.Connection.Credentials = []byte(*credentials)
	}

	return fields, nil
}
