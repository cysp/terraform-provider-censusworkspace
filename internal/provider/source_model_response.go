package provider

import (
	"context"
	"strconv"

	cm "github.com/cysp/terraform-provider-censusworkspace/internal/census-management-go"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewSourceModelFromResponse(_ context.Context, response cm.SourceData) (SourceModel, diag.Diagnostics) {
	model := SourceModel{
		sourceModelBase: sourceModelBase{
			ID:    types.StringValue(strconv.FormatInt(response.ID, 10)),
			Name:  types.StringValue(response.Name),
			Label: types.StringPointerValue(response.Label.ValueStringPointer()),
		},
		Type: types.StringValue(response.Type),
	}

	if syncEngine, syncEngineOk := response.SyncEngine.Get(); syncEngineOk {
		model.SyncEngine = types.StringValue(syncEngine)
	}

	if warehouseWritebackRetentionInDays, warehouseWritebackRetentionInDaysOk := response.WarehouseWritebackRetentionInDays.Get(); warehouseWritebackRetentionInDaysOk {
		model.WarehouseWritebackRetentionInDays = types.Int64Value(warehouseWritebackRetentionInDays)
	}

	if response.ConnectionDetails != nil {
		model.ConnectionDetails = jsontypes.NewNormalizedValue(string(response.ConnectionDetails))
	}

	return model, nil
}
