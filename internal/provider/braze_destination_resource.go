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
	_ resource.Resource                = (*brazeDestinationResource)(nil)
	_ resource.ResourceWithConfigure   = (*brazeDestinationResource)(nil)
	_ resource.ResourceWithIdentity    = (*brazeDestinationResource)(nil)
	_ resource.ResourceWithImportState = (*brazeDestinationResource)(nil)
	_ resource.ResourceWithMoveState   = (*brazeDestinationResource)(nil)
)

//nolint:ireturn
func NewBrazeDestinationResource() resource.Resource {
	return &brazeDestinationResource{}
}

type brazeDestinationResource struct {
	providerData ProviderData
}

func (r *brazeDestinationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_braze_destination"
}

func (r *brazeDestinationResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = BrazeDestinationResourceSchema(ctx)
}

func (r *brazeDestinationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	resp.Diagnostics.Append(SetProviderDataFromResourceConfigureRequest(req, &r.providerData)...)
}

func (r *brazeDestinationResource) IdentitySchema(ctx context.Context, _ resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = BrazeDestinationResourceIdentitySchema(ctx)
}

func (r *brazeDestinationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	ImportStatePassthroughMultipartID(ctx, []path.Path{path.Root("id")}, req, resp)
}

func (r *brazeDestinationResource) MoveState(ctx context.Context) []resource.StateMover {
	schema := DestinationResourceSchema(ctx)

	return []resource.StateMover{
		{
			SourceSchema: &schema,
			StateMover: func(_ context.Context, req resource.MoveStateRequest, resp *resource.MoveStateResponse) {
				if req.SourceTypeName == DestinationResourceTypeName && req.SourceSchemaVersion == 0 {
					destinationModel := DestinationModel{}
					resp.Diagnostics.Append(req.SourceState.Get(ctx, &destinationModel)...)

					if destinationModel.Type.ValueString() != BrazeDestinationType {
						return
					}

					type destinationCredentialsModel struct {
						InstanceURL *string `json:"instance_url"`
						APIKey      *string `json:"api_key"`
						ClientKey   *string `json:"client_key"`
					}

					destinationCredentials := destinationCredentialsModel{}
					resp.Diagnostics.Append(destinationModel.Credentials.Unmarshal(&destinationCredentials)...)

					brazeDestinationCredentials := BrazeDestinationCredentials{
						InstanceURL: types.StringPointerValue(destinationCredentials.InstanceURL),
						APIKey:      types.StringPointerValue(destinationCredentials.APIKey),
						ClientKey:   types.StringPointerValue(destinationCredentials.ClientKey),
					}

					type destinationConnectionDetailsModel struct {
						InstanceURL *string `json:"instance_url"`
					}

					destinationConnectionDetails := destinationConnectionDetailsModel{}
					resp.Diagnostics.Append(destinationModel.ConnectionDetails.Unmarshal(&destinationConnectionDetails)...)

					brazeDestinationConnectionDetails := BrazeDestinationConnectionDetails{
						InstanceURL: types.StringPointerValue(destinationConnectionDetails.InstanceURL),
					}

					brazeDestinationModel := BrazeDestinationModel{
						destinationModelBase: destinationModel.destinationModelBase,
						Credentials:          NewTypedObject(brazeDestinationCredentials),
						ConnectionDetails:    NewTypedObject(brazeDestinationConnectionDetails),
					}

					resp.Diagnostics.Append(resp.TargetState.Set(ctx, &brazeDestinationModel)...)

					return
				}
			},
		},
	}
}

//nolint:dupl
func (r *brazeDestinationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan BrazeDestinationModel

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

	model, modelDiags := NewBrazeDestinationModelFromResponse(ctx, getDestinationResponse.Response.Data)
	resp.Diagnostics.Append(modelDiags...)

	credentials := plan.Credentials.Value()
	credentials.UpdateWithConnectionDetails(model.ConnectionDetails.Value())

	model.Credentials = NewTypedObject(credentials)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.Identity.SetAttribute(ctx, path.Root("id"), model.ID)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func (r *brazeDestinationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state BrazeDestinationModel

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

	model, modelDiags := NewBrazeDestinationModelFromResponse(ctx, getDestinationResponse.Response.Data)
	resp.Diagnostics.Append(modelDiags...)

	credentials := state.Credentials.Value()
	credentials.UpdateWithConnectionDetails(model.ConnectionDetails.Value())

	model.Credentials = NewTypedObject(credentials)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.Identity.SetAttribute(ctx, path.Root("id"), model.ID)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func (r *brazeDestinationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state, plan BrazeDestinationModel

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

	model, modelDiags := NewBrazeDestinationModelFromResponse(ctx, updateDestinationResponse.Response.Data)
	resp.Diagnostics.Append(modelDiags...)

	credentials := plan.Credentials.Value()
	credentials.UpdateWithConnectionDetails(model.ConnectionDetails.Value())

	model.Credentials = NewTypedObject(credentials)

	resp.Diagnostics.Append(resp.Identity.SetAttribute(ctx, path.Root("id"), model.ID)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

//nolint:dupl
func (r *brazeDestinationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state BrazeDestinationModel

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
