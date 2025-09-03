package provider

import (
	"context"

	cm "github.com/cysp/terraform-provider-censusworkspace/internal/census-management-go"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func (m *SQLDatasetModel) ToCreateDatasetBody(ctx context.Context) (cm.CreateDatasetBody, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	body := cm.NewCreateSQLDatasetBodyCreateDatasetBody(m.ToCreateSQLDatasetBody(ctx))

	return body, diags
}

func (m *SQLDatasetModel) ToCreateSQLDatasetBody(_ context.Context) cm.CreateSQLDatasetBody {
	body := cm.CreateSQLDatasetBody{
		Name:     m.Name.ValueString(),
		Type:     cm.CreateSQLDatasetBodyTypeSQL,
		SourceID: m.SourceID.ValueInt64(),
		Query:    m.Query.ValueString(),
	}

	description := m.Description.ValueString()
	if description != "" {
		body.Description.SetTo(description)
	} else {
		body.Description.SetToNull()
	}

	return body
}

func (m *SQLDatasetModel) ToUpdateDatasetBody(ctx context.Context) (cm.UpdateDatasetBody, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	body := cm.NewUpdateSQLDatasetBodyUpdateDatasetBody(m.ToUpdateSQLDatasetBody(ctx))

	return body, diags
}

func (m *SQLDatasetModel) ToUpdateSQLDatasetBody(_ context.Context) cm.UpdateSQLDatasetBody {
	body := cm.UpdateSQLDatasetBody{
		Name:  cm.NewOptString(m.Name.ValueString()),
		Query: cm.NewOptString(m.Query.ValueString()),
	}

	description := m.Description.ValueString()
	if description != "" {
		body.Description.SetTo(description)
	} else {
		body.Description.SetToNull()
	}

	return body
}
