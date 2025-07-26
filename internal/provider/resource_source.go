package provider

import (
	"context"

	cm "github.com/cysp/terraform-provider-censusworkspace/internal/census-management-go/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
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
	providerData CensusProviderData
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

func (r *sourceResource) IdentitySchema(_ context.Context, _ resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"id": identityschema.StringAttribute{RequiredForImport: true},
		},
	}
}

func (r *sourceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughWithIdentity(ctx, path.Root("id"), path.Root("id"), req, resp)
}

func (r *sourceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data SourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	request, requestDiags := data.ToCreateSourceData(ctx)
	resp.Diagnostics.Append(requestDiags...)

	if resp.Diagnostics.HasError() {
		return
	}

	response, err := r.providerData.client.CreateSource(ctx, &request)

	tflog.Info(ctx, "source.create", map[string]interface{}{
		"request":  request,
		"response": response,
		"err":      err,
	})

	if response == nil {
		resp.Diagnostics.AddError("Failed to create source", err.Error())
		return
	}

	responseModel, responseModelDiags := NewSourceResourceModelFromResponse(ctx, response.Response.Data)
	resp.Diagnostics.Append(responseModelDiags...)

	data = responseModel

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *sourceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data SourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	params := cm.GetSourceParams{
		SourceID: data.ID.ValueInt64(),
	}

	response, err := r.providerData.client.GetSource(ctx, params)

	tflog.Info(ctx, "source.read", map[string]interface{}{
		"params":   params,
		"response": response,
		"err":      err,
	})

	if response == nil {
		resp.Diagnostics.AddError("Failed to read source", "")
		return
	}

	responseModel, responseModelDiags := NewSourceResourceModelFromResponse(ctx, response.Response.Data)
	resp.Diagnostics.Append(responseModelDiags...)

	data = responseModel

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *sourceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data SourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	params := cm.UpdateSourceParams{
		SourceID: data.ID.ValueInt64(),
	}

	request, requestDiags := data.ToUpdateSourceData(ctx)
	resp.Diagnostics.Append(requestDiags...)

	if resp.Diagnostics.HasError() {
		return
	}

	response, err := r.providerData.client.UpdateSource(ctx, &request, params)

	tflog.Info(ctx, "source.update", map[string]interface{}{
		"params":   params,
		"request":  request,
		"response": response,
		"err":      err,
	})

	if response == nil {
		resp.Diagnostics.AddError("Failed to update source", "")
		return
	}

	responseModel, responseModelDiags := NewSourceResourceModelFromResponse(ctx, response.Response.Data)
	resp.Diagnostics.Append(responseModelDiags...)

	data = responseModel

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *sourceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data SourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	params := cm.DeleteSourceParams{
		SourceID: data.ID.ValueInt64(),
	}

	response, err := r.providerData.client.DeleteSource(ctx, params)

	tflog.Info(ctx, "source.delete", map[string]interface{}{
		"params":   params,
		"response": response,
		"err":      err,
	})

	if response == nil || response.Response.Status != cm.ResponseStatusSuccess {
		resp.Diagnostics.AddError("Failed to delete source", response.Response.Message.Or(err.Error()))
		return
	}
}
