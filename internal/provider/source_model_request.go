package provider

import (
	"context"

	cm "github.com/cysp/terraform-provider-censusworkspace/internal/census-management-go/client"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func (m *SourceModel) ToCreateSourceData(_ context.Context) (cm.CreateSourceData, diag.Diagnostics) {
	fields := cm.CreateSourceData{
		Type:  m.Type.ValueString(),
		Label: cm.NewNilPointerString(m.Type.ValueStringPointer()),
	}

	if credentials := m.Credentials.ValueStringPointer(); credentials != nil {
		fields.Credentials = []byte(*credentials)
	}

	return fields, nil
}

func (m *SourceModel) ToUpdateSourceData(_ context.Context) (cm.UpdateSourceData, diag.Diagnostics) {
	fields := cm.UpdateSourceData{
		Label: cm.NewNilPointerString(m.Type.ValueStringPointer()),
	}

	if credentials := m.Credentials.ValueStringPointer(); credentials != nil {
		fields.Credentials = []byte(*credentials)
	}

	return fields, nil
}
