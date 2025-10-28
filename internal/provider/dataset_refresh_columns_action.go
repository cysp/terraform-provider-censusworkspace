package provider

import (
	"context"
	"strconv"

	cm "github.com/cysp/terraform-provider-censusworkspace/internal/census-management-go"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ action.Action              = (*datasetRefreshColumnsAction)(nil)
	_ action.ActionWithConfigure = (*datasetRefreshColumnsAction)(nil)
)

//nolint:ireturn
func NewDatasetRefreshColumnsAction() action.Action {
	return &datasetRefreshColumnsAction{}
}

type datasetRefreshColumnsAction struct {
	providerData ProviderData
}

func (r *datasetRefreshColumnsAction) Metadata(_ context.Context, req action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dataset_refresh_columns"
}

func (r *datasetRefreshColumnsAction) Schema(ctx context.Context, _ action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = DatasetRefreshColumnsActionSchema(ctx)
}

func (r *datasetRefreshColumnsAction) Configure(_ context.Context, req action.ConfigureRequest, resp *action.ConfigureResponse) {
	resp.Diagnostics.Append(SetProviderDataFromActionConfigureRequest(req, &r.providerData)...)
}

// Invoke implements action.ActionWithConfigure.
func (r *datasetRefreshColumnsAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var config DatasetRefreshColumnsActionModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	refreshDatasetColumnsParams := cm.RefreshDatasetColumnsParams{
		DatasetID: strconv.FormatInt(config.DatasetID.ValueInt64(), 10),
	}

	refreshDatasetColumnsResponse, refreshDatasetColumnsErr := r.providerData.client.RefreshDatasetColumns(ctx, refreshDatasetColumnsParams)

	tflog.Info(ctx, "dataset_refresh_columns.invoke", map[string]any{
		"params":   refreshDatasetColumnsParams,
		"response": refreshDatasetColumnsResponse,
		"err":      refreshDatasetColumnsErr,
	})

	if refreshDatasetColumnsResponse == nil {
		resp.Diagnostics.AddError("Failed to refresh dataset columns", refreshDatasetColumnsErr.Error())

		return
	}

	resp.SendProgress(action.InvokeProgressEvent{
		Message: "Dataset columns refresh initiated: " + strconv.FormatInt(refreshDatasetColumnsResponse.Response.RefreshKey, 10),
	})
}
