package provider

import (
	"context"
	"strconv"

	cm "github.com/cysp/terraform-provider-censusworkspace/internal/census-management-go"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewSQLDatasetModelFromResponse(_ context.Context, path path.Path, data cm.DatasetData) (SQLDatasetModel, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	model := SQLDatasetModel{}

	sql, ok := data.GetSQLDatasetData()
	if !ok {
		diags.AddAttributeError(path.AtName("type"), "Incorrect dataset type", "Expected sql dataset, got "+string(data.Type))

		return model, diags
	}

	model.ID = types.StringValue(strconv.FormatInt(sql.ID, 10))
	model.Name = types.StringValue(sql.Name)
	model.SourceID = types.Int64Value(sql.SourceID)

	model.Query = types.StringValue(sql.Query)

	description, _ := sql.Description.Get()
	if description != "" {
		model.Description = types.StringValue(description)
	}

	model.CreatedAt = timetypes.NewRFC3339TimeValue(sql.CreatedAt)
	model.UpdatedAt = timetypes.NewRFC3339TimeValue(sql.UpdatedAt)

	return model, diags
}
