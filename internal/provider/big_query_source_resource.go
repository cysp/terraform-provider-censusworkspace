//nolint:dupl
package provider

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	cm "github.com/cysp/terraform-provider-censusworkspace/internal/census-management-go"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = (*bigQuerySourceResource)(nil)
	_ resource.ResourceWithConfigure   = (*bigQuerySourceResource)(nil)
	_ resource.ResourceWithIdentity    = (*bigQuerySourceResource)(nil)
	_ resource.ResourceWithImportState = (*bigQuerySourceResource)(nil)
	_ resource.ResourceWithMoveState   = (*bigQuerySourceResource)(nil)
)

//nolint:ireturn
func NewBigQuerySourceResource() resource.Resource {
	return &bigQuerySourceResource{}
}

type bigQuerySourceResource struct {
	providerData ProviderData
}

func (r *bigQuerySourceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_big_query_source"
}

func (r *bigQuerySourceResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = BigQuerySourceResourceSchema(ctx)
}

func (r *bigQuerySourceResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	resp.Diagnostics.Append(SetProviderDataFromResourceConfigureRequest(req, &r.providerData)...)
}

func (r *bigQuerySourceResource) IdentitySchema(ctx context.Context, _ resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = BigQuerySourceResourceIdentitySchema(ctx)
}

func (r *bigQuerySourceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughWithIdentity(ctx, path.Root("id"), path.Root("id"), req, resp)
}

func (r *bigQuerySourceResource) MoveState(ctx context.Context) []resource.StateMover {
	schema := SourceResourceSchema(ctx)

	return []resource.StateMover{
		{
			SourceSchema: &schema,
			StateMover: func(_ context.Context, req resource.MoveStateRequest, resp *resource.MoveStateResponse) {
				if req.SourceTypeName == "censusworkspace_source" && req.SourceSchemaVersion == 0 {
					sourceModel := SourceModel{}
					resp.Diagnostics.Append(req.SourceState.Get(ctx, &sourceModel)...)

					if sourceModel.Type.ValueString() != BigQuerySourceType {
						return
					}

					type sourceCredentialsModel struct {
						ProjectID         *string `json:"project_id"`
						Location          *string `json:"location"`
						ServiceAccountKey *string `json:"service_account_key"`
					}

					sourceCredentials := sourceCredentialsModel{}
					resp.Diagnostics.Append(sourceModel.Credentials.Unmarshal(&sourceCredentials)...)

					bigQuerySourceCredentials := BigQuerySourceCredentials{
						ProjectID:         types.StringPointerValue(sourceCredentials.ProjectID),
						Location:          types.StringPointerValue(sourceCredentials.Location),
						ServiceAccountKey: types.StringPointerValue(sourceCredentials.ServiceAccountKey),
					}

					type sourceConnectionDetailsModel struct {
						ProjectID      *string `json:"project_id"`
						Location       *string `json:"location"`
						ServiceAccount *string `json:"service_account"`
					}

					sourceConnectionDetails := sourceConnectionDetailsModel{}
					resp.Diagnostics.Append(sourceModel.ConnectionDetails.Unmarshal(&sourceConnectionDetails)...)

					bigQuerySourceConnectionDetails := BigQuerySourceConnectionDetails{
						ProjectID:      types.StringPointerValue(sourceConnectionDetails.ProjectID),
						Location:       types.StringPointerValue(sourceConnectionDetails.Location),
						ServiceAccount: types.StringPointerValue(sourceConnectionDetails.ServiceAccount),
					}

					bigQuerySourceModel := BigQuerySourceModel{
						sourceModelBase:   sourceModel.sourceModelBase,
						Credentials:       NewTypedObject(bigQuerySourceCredentials),
						ConnectionDetails: NewTypedObject(bigQuerySourceConnectionDetails),
					}

					resp.Diagnostics.Append(resp.TargetState.Set(ctx, &bigQuerySourceModel)...)

					return
				}
			},
		},
	}
}

func (r *bigQuerySourceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan BigQuerySourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	createSourceRequest, createSourceRequestDiags := plan.ToCreateSourceData(ctx)
	resp.Diagnostics.Append(createSourceRequestDiags...)

	if resp.Diagnostics.HasError() {
		return
	}

	createSourceResponse, createSourceErr := r.providerData.client.CreateSource(ctx, &createSourceRequest)

	tflog.Info(ctx, "source.create", map[string]any{
		"request":  createSourceRequest,
		"response": createSourceResponse,
		"err":      createSourceErr,
	})

	if createSourceResponse == nil {
		resp.Diagnostics.AddError("Failed to create source", createSourceErr.Error())

		return
	}

	sourceID := createSourceResponse.Response.Data.ID
	sourceIDString := strconv.FormatInt(sourceID, 10)

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), sourceIDString)...)

	if resp.Diagnostics.HasError() {
		return
	}

	getSourceParams := cm.GetSourceParams{
		SourceID: sourceIDString,
	}

	getSourceResponse, getSourceErr := r.providerData.client.GetSource(ctx, getSourceParams)

	tflog.Info(ctx, "source.read", map[string]any{
		"request":  getSourceParams,
		"response": getSourceResponse,
		"err":      getSourceErr,
	})

	model, modelDiags := NewBigQuerySourceModelFromResponse(ctx, getSourceResponse.Response.Data)
	resp.Diagnostics.Append(modelDiags...)

	if model.SyncEngine.IsNull() && !plan.SyncEngine.IsUnknown() {
		model.SyncEngine = plan.SyncEngine
	}

	model.Credentials = plan.Credentials

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func (r *bigQuerySourceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state BigQuerySourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	params := cm.GetSourceParams{
		SourceID: state.ID.ValueString(),
	}

	getSourceResponse, getSourceErr := r.providerData.client.GetSource(ctx, params)

	tflog.Info(ctx, "source.read", map[string]any{
		"params":   params,
		"response": getSourceResponse,
		"err":      getSourceErr,
	})

	if getSourceResponse == nil {
		var srsc *cm.StatusResponseStatusCode
		if errors.As(getSourceErr, &srsc) {
			if srsc.StatusCode == http.StatusNotFound {
				resp.Diagnostics.AddWarning("Failed to read source", srsc.Error())
				resp.State.RemoveResource(ctx)

				return
			}
		}

		resp.Diagnostics.AddError("Failed to read source", getSourceErr.Error())

		return
	}

	model, modelDiags := NewBigQuerySourceModelFromResponse(ctx, getSourceResponse.Response.Data)
	resp.Diagnostics.Append(modelDiags...)

	if model.SyncEngine.IsNull() && !state.SyncEngine.IsUnknown() {
		model.SyncEngine = state.SyncEngine
	}

	model.Credentials = state.Credentials

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func (r *bigQuerySourceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state, plan BigQuerySourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	params := cm.UpdateSourceParams{
		SourceID: plan.ID.ValueString(),
	}

	updateSourceRequest, updateSourceRequestDiags := plan.ToUpdateSourceData(ctx)
	resp.Diagnostics.Append(updateSourceRequestDiags...)

	if resp.Diagnostics.HasError() {
		return
	}

	updateSourceResponse, updateSourceErr := r.providerData.client.UpdateSource(ctx, &updateSourceRequest, params)

	tflog.Info(ctx, "source.update", map[string]any{
		"params":   params,
		"request":  updateSourceRequest,
		"response": updateSourceResponse,
		"err":      updateSourceErr,
	})

	if updateSourceResponse == nil {
		resp.Diagnostics.AddError("Failed to update source", updateSourceErr.Error())

		return
	}

	model, modelDiags := NewBigQuerySourceModelFromResponse(ctx, updateSourceResponse.Response.Data)
	resp.Diagnostics.Append(modelDiags...)

	if model.SyncEngine.IsNull() && !state.SyncEngine.IsUnknown() {
		model.SyncEngine = state.SyncEngine
	}

	model.Credentials = state.Credentials

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func (r *bigQuerySourceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state BigQuerySourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	params := cm.DeleteSourceParams{
		SourceID: state.ID.ValueString(),
	}

	deleteSourceResponse, deleteSourceErr := r.providerData.client.DeleteSource(ctx, params)

	tflog.Info(ctx, "source.delete", map[string]any{
		"params":   params,
		"response": deleteSourceResponse,
		"err":      deleteSourceErr,
	})

	var srsc *cm.StatusResponseStatusCode
	if errors.As(deleteSourceErr, &srsc) {
		if srsc.StatusCode == http.StatusNotFound {
			resp.Diagnostics.AddWarning("Source not found", srsc.Error())
			resp.State.RemoveResource(ctx)

			return
		}
	}

	if deleteSourceResponse == nil || deleteSourceResponse.Response.Status.ResponseStatus != cm.ResponseStatusDeleted {
		var detail string

		if deleteSourceResponse != nil {
			detail = deleteSourceResponse.Response.Message.Value
		}

		if detail == "" && deleteSourceErr != nil {
			detail = deleteSourceErr.Error()
		}

		resp.Diagnostics.AddError("Failed to delete source", detail)

		return
	}
}
