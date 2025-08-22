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
	_ resource.Resource                = (*destinationResource)(nil)
	_ resource.ResourceWithConfigure   = (*destinationResource)(nil)
	_ resource.ResourceWithIdentity    = (*destinationResource)(nil)
	_ resource.ResourceWithImportState = (*destinationResource)(nil)
)

//nolint:ireturn
func NewDestinationResource() resource.Resource {
	return &destinationResource{}
}

type destinationResource struct {
	providerData ProviderData
}

func (r *destinationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_destination"
}

func (r *destinationResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = DestinationResourceSchema(ctx)
}

func (r *destinationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	resp.Diagnostics.Append(SetProviderDataFromResourceConfigureRequest(req, &r.providerData)...)
}

func (r *destinationResource) IdentitySchema(ctx context.Context, _ resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = DestinationResourceIdentitySchema(ctx)
}

func (r *destinationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughWithIdentity(ctx, path.Root("id"), path.Root("id"), req, resp)
}

func (r *destinationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan DestinationModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	createDestinationRequest, createDestinationRequestDiags := plan.ToCreateDestinationData(ctx)
	resp.Diagnostics.Append(createDestinationRequestDiags...)

	if resp.Diagnostics.HasError() {
		return
	}

	createDestinationResponse, createDestinationErr := r.providerData.client.CreateDestination(ctx, &createDestinationRequest)

	tflog.Info(ctx, "destination.create", map[string]any{
		"request":  createDestinationRequest,
		"response": createDestinationResponse,
		"err":      createDestinationErr,
	})

	if createDestinationResponse == nil {
		resp.Diagnostics.AddError("Failed to create destination", createDestinationErr.Error())

		return
	}

	destinationID := createDestinationResponse.Response.Data.ID
	destinationIDString := strconv.FormatInt(destinationID, 10)

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), destinationIDString)...)

	if resp.Diagnostics.HasError() {
		return
	}

	getDestinationParams := cm.GetDestinationParams{
		DestinationID: destinationIDString,
	}

	getDestinationResponse, getDestinationErr := r.providerData.client.GetDestination(ctx, getDestinationParams)

	tflog.Info(ctx, "destination.read", map[string]any{
		"request":  getDestinationParams,
		"response": getDestinationResponse,
		"err":      getDestinationErr,
	})

	model, modelDiags := NewDestinationModelFromResponse(ctx, getDestinationResponse.Response.Data)
	resp.Diagnostics.Append(modelDiags...)

	model.Credentials = plan.Credentials

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func (r *destinationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state DestinationModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	params := cm.GetDestinationParams{
		DestinationID: state.ID.ValueString(),
	}

	getDestinationResponse, getDestinationErr := r.providerData.client.GetDestination(ctx, params)

	tflog.Info(ctx, "destination.read", map[string]any{
		"params":   params,
		"response": getDestinationResponse,
		"err":      getDestinationErr,
	})

	if getDestinationResponse == nil {
		var srsc *cm.StatusResponseStatusCode
		if errors.As(getDestinationErr, &srsc) {
			if srsc.StatusCode == http.StatusNotFound {
				resp.Diagnostics.AddWarning("Failed to read destination", srsc.Error())
				resp.State.RemoveResource(ctx)

				return
			}
		}

		resp.Diagnostics.AddError("Failed to read destination", getDestinationErr.Error())

		return
	}

	model, modelDiags := NewDestinationModelFromResponse(ctx, getDestinationResponse.Response.Data)
	resp.Diagnostics.Append(modelDiags...)

	model.Credentials = state.Credentials

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func (r *destinationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state, plan DestinationModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	params := cm.UpdateDestinationParams{
		DestinationID: plan.ID.ValueString(),
	}

	updateDestinationRequest, updateDestinationRequestDiags := plan.ToUpdateDestinationData(ctx)
	resp.Diagnostics.Append(updateDestinationRequestDiags...)

	if resp.Diagnostics.HasError() {
		return
	}

	updateDestinationResponse, updateDestinationErr := r.providerData.client.UpdateDestination(ctx, &updateDestinationRequest, params)

	tflog.Info(ctx, "destination.update", map[string]any{
		"params":   params,
		"request":  updateDestinationRequest,
		"response": updateDestinationResponse,
		"err":      updateDestinationErr,
	})

	if updateDestinationResponse == nil {
		resp.Diagnostics.AddError("Failed to update destination", updateDestinationErr.Error())

		return
	}

	model, modelDiags := NewDestinationModelFromResponse(ctx, updateDestinationResponse.Response.Data)
	resp.Diagnostics.Append(modelDiags...)

	model.Credentials = state.Credentials

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

//nolint:dupl
func (r *destinationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state DestinationModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	params := cm.DeleteDestinationParams{
		DestinationID: state.ID.ValueString(),
	}

	deleteDestinationResponse, deleteDestinationErr := r.providerData.client.DeleteDestination(ctx, params)

	tflog.Info(ctx, "destination.delete", map[string]any{
		"params":   params,
		"response": deleteDestinationResponse,
		"err":      deleteDestinationErr,
	})

	var srsc *cm.StatusResponseStatusCode
	if errors.As(deleteDestinationErr, &srsc) {
		if srsc.StatusCode == http.StatusNotFound {
			resp.Diagnostics.AddWarning("Destination not found", srsc.Error())
			resp.State.RemoveResource(ctx)

			return
		}
	}

	if deleteDestinationResponse == nil || deleteDestinationResponse.Response.Status.ResponseStatus != cm.ResponseStatusDeleted {
		var detail string

		if deleteDestinationResponse != nil {
			detail = deleteDestinationResponse.Response.Message.Value
		}

		if detail == "" && deleteDestinationErr != nil {
			detail = deleteDestinationErr.Error()
		}

		resp.Diagnostics.AddError("Failed to delete destination", detail)

		return
	}
}
