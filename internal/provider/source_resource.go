package provider

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	cm "github.com/cysp/terraform-provider-censusworkspace/internal/census-management-go"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = (*sourceResource)(nil)
	_ resource.ResourceWithConfigure   = (*sourceResource)(nil)
	_ resource.ResourceWithIdentity    = (*sourceResource)(nil)
	_ resource.ResourceWithImportState = (*sourceResource)(nil)
)

//nolint:ireturn
func NewSourceResource() resource.Resource {
	return &sourceResource{}
}

type sourceResource struct {
	providerData ProviderData
}

func (r *sourceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_source"
}

func (r *sourceResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = SourceResourceSchema(ctx)
}

func (r *sourceResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	resp.Diagnostics.Append(SetProviderDataFromResourceConfigureRequest(req, &r.providerData)...)
}

func (r *sourceResource) IdentitySchema(ctx context.Context, _ resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = SourceResourceIdentitySchema(ctx)
}

func (r *sourceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughWithIdentity(ctx, path.Root("id"), path.Root("id"), req, resp)
}

func (r *sourceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan SourceModel

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

	model, modelDiags := NewSourceModelFromResponse(ctx, getSourceResponse.Response.Data)
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

func (r *sourceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state SourceModel

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

	model, modelDiags := NewSourceModelFromResponse(ctx, getSourceResponse.Response.Data)
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

func (r *sourceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state, plan SourceModel

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

	model, modelDiags := NewSourceModelFromResponse(ctx, updateSourceResponse.Response.Data)
	resp.Diagnostics.Append(modelDiags...)

	if model.SyncEngine.IsNull() && !state.SyncEngine.IsUnknown() {
		model.SyncEngine = state.SyncEngine
	}

	model.Credentials = state.Credentials

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func (r *sourceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state SourceModel

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
