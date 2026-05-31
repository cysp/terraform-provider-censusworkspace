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
	_ resource.Resource                   = (*customAPIDestinationResource)(nil)
	_ resource.ResourceWithConfigure      = (*customAPIDestinationResource)(nil)
	_ resource.ResourceWithIdentity       = (*customAPIDestinationResource)(nil)
	_ resource.ResourceWithImportState    = (*customAPIDestinationResource)(nil)
	_ resource.ResourceWithModifyPlan     = (*customAPIDestinationResource)(nil)
	_ resource.ResourceWithMoveState      = (*customAPIDestinationResource)(nil)
	_ resource.ResourceWithValidateConfig = (*customAPIDestinationResource)(nil)
)

//nolint:ireturn
func NewCustomAPIDestinationResource() resource.Resource {
	return &customAPIDestinationResource{}
}

type customAPIDestinationResource struct {
	providerData ProviderData
}

func (r *customAPIDestinationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_custom_api_destination"
}

func (r *customAPIDestinationResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = CustomAPIDestinationResourceSchema(ctx)
}

func (r *customAPIDestinationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	resp.Diagnostics.Append(SetProviderDataFromResourceConfigureRequest(req, &r.providerData)...)
}

func (r *customAPIDestinationResource) IdentitySchema(ctx context.Context, _ resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = CustomAPIDestinationResourceIdentitySchema(ctx)
}

func (r *customAPIDestinationResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var config CustomAPIDestinationModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	resp.Diagnostics.Append(hydrateCustomAPIHeaderWriteOnlyValues(ctx, req.Config, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	credentials := config.Credentials.Value()
	if credentials.CustomHeaders.IsNull() || credentials.CustomHeaders.IsUnknown() {
		return
	}

	for key, header := range credentials.CustomHeaders.Elements() {
		headerValue := header.Value()
		validateStringCredential(&resp.Diagnostics, headerValue.Value, headerValue.ValueWO, customAPIHeaderValuePath(key), customAPIHeaderValueWOPath(key))
	}
}

func (r *customAPIDestinationResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() {
		return
	}

	var plan, config CustomAPIDestinationModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	resp.Diagnostics.Append(hydrateCustomAPIHeaderWriteOnlyValues(ctx, req.Config, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, values, modelDiags := customAPIDestinationModelWithWriteOnlyCredentials(plan, config)
	resp.Diagnostics.Append(modelDiags...)

	markWriteOnlyCredentialChange(ctx, req, resp, values, NewTypedObjectUnknown[CustomAPIDestinationConnectionDetails]())
}

func (r *customAPIDestinationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	ImportStatePassthroughMultipartID(ctx, []path.Path{path.Root("id")}, req, resp)
}

func (r *customAPIDestinationResource) MoveState(ctx context.Context) []resource.StateMover {
	schema := DestinationResourceSchema(ctx)

	return []resource.StateMover{
		{
			SourceSchema: &schema,
			StateMover: func(_ context.Context, req resource.MoveStateRequest, resp *resource.MoveStateResponse) {
				if req.SourceTypeName == DestinationResourceTypeName && req.SourceSchemaVersion == 0 {
					destinationModel := DestinationModel{}
					resp.Diagnostics.Append(req.SourceState.Get(ctx, &destinationModel)...)

					if destinationModel.Type.ValueString() != CustomAPIDestinationType {
						return
					}

					type destinationCredentialsCustomHeaderModel struct {
						Value    *string `json:"value"`
						IsSecret *bool   `json:"is_secret"`
					}

					type destinationCredentialsModel struct {
						APIVersion    *int64                                             `json:"api_version"`
						WebhookURL    *string                                            `json:"webhook_url"`
						CustomHeaders map[string]destinationCredentialsCustomHeaderModel `json:"custom_headers"`
					}

					destinationCredentials := destinationCredentialsModel{}
					resp.Diagnostics.Append(destinationModel.Credentials.Unmarshal(&destinationCredentials)...)

					customAPIDestinationCredentials := CustomAPIDestinationCredentials{
						APIVersion: types.Int64PointerValue(destinationCredentials.APIVersion),
						WebhookURL: types.StringPointerValue(destinationCredentials.WebhookURL),
					}

					if destinationCredentials.CustomHeaders != nil {
						customAPIDestinationCredentialsCustomHeaders := make(map[string]TypedObject[CustomAPIDestinationCustomHeader], 0)
						for key, value := range destinationCredentials.CustomHeaders {
							customAPIDestinationCredentialsCustomHeaders[key] = NewTypedObject(CustomAPIDestinationCustomHeader{
								Value:    types.StringPointerValue(value.Value),
								IsSecret: types.BoolPointerValue(value.IsSecret),
							})
						}

						customAPIDestinationCredentials.CustomHeaders = NewTypedMap(customAPIDestinationCredentialsCustomHeaders)
					}

					type destinationConnectionDetailsCustomHeaderModel struct {
						Value    *string `json:"value"`
						IsSecret *bool   `json:"is_secret"`
					}

					type destinationConnectionDetailsModel struct {
						APIVersion    *int64                                                   `json:"api_version"`
						WebhookURL    *string                                                  `json:"webhook_url"`
						CustomHeaders map[string]destinationConnectionDetailsCustomHeaderModel `json:"custom_headers"`
					}

					destinationConnectionDetails := destinationConnectionDetailsModel{}
					resp.Diagnostics.Append(destinationModel.ConnectionDetails.Unmarshal(&destinationConnectionDetails)...)

					customAPIDestinationConnectionDetails := CustomAPIDestinationConnectionDetails{
						APIVersion: types.Int64PointerValue(destinationConnectionDetails.APIVersion),
						WebhookURL: types.StringPointerValue(destinationConnectionDetails.WebhookURL),
					}

					if destinationConnectionDetails.CustomHeaders != nil {
						customAPIDestinationConnectionDetailsCustomHeaders := make(map[string]TypedObject[CustomAPIDestinationConnectionDetailsCustomHeader], 0)
						for key, value := range destinationConnectionDetails.CustomHeaders {
							customAPIDestinationConnectionDetailsCustomHeaders[key] = NewTypedObject(CustomAPIDestinationConnectionDetailsCustomHeader{
								Value:    types.StringPointerValue(value.Value),
								IsSecret: types.BoolPointerValue(value.IsSecret),
							})
						}

						customAPIDestinationConnectionDetails.CustomHeaders = NewTypedMap(customAPIDestinationConnectionDetailsCustomHeaders)
					}

					customAPIDestinationModel := CustomAPIDestinationModel{
						destinationModelBase: destinationModel.destinationModelBase,
						Credentials:          NewTypedObject(customAPIDestinationCredentials),
						ConnectionDetails:    NewTypedObject(customAPIDestinationConnectionDetails),
					}

					resp.Diagnostics.Append(resp.TargetState.Set(ctx, &customAPIDestinationModel)...)

					return
				}
			},
		},
	}
}

func (r *customAPIDestinationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, config CustomAPIDestinationModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	resp.Diagnostics.Append(hydrateCustomAPIHeaderWriteOnlyValues(ctx, req.Config, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	requestModel, writeOnlyValues, requestModelDiags := customAPIDestinationModelWithWriteOnlyCredentials(plan, config)
	resp.Diagnostics.Append(requestModelDiags...)

	if resp.Diagnostics.HasError() {
		return
	}

	createDestinationRequest, createDestinationRequestDiags := requestModel.ToCreateDestinationData(ctx)
	resp.Diagnostics.Append(createDestinationRequestDiags...)

	if resp.Diagnostics.HasError() {
		return
	}

	createDestinationResponse, createDestinationErr := r.providerData.client.CreateDestination(ctx, &createDestinationRequest)

	tflog.Info(ctx, "destination.create", map[string]any{
		"response": createDestinationResponse,
		"err":      createDestinationErr,
	})

	if createDestinationResponse == nil {
		resp.Diagnostics.AddError("Failed to create destination", detailFromError(createDestinationErr))

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

	if getDestinationResponse == nil {
		resp.Diagnostics.AddError("Failed to read destination after create", detailFromError(getDestinationErr))

		return
	}

	model, modelDiags := NewCustomAPIDestinationModelFromResponse(ctx, getDestinationResponse.Response.Data)
	resp.Diagnostics.Append(modelDiags...)

	model = sanitizedCustomAPIDestinationCredentials(model, plan, config, model.ConnectionDetails.Value(), writeOnlyValues)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.Identity.SetAttribute(ctx, path.Root("id"), model.ID)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
	resp.Diagnostics.Append(writeWriteOnlyCredentialVerifiers(ctx, resp.Private, writeOnlyValues)...)
}

func (r *customAPIDestinationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state CustomAPIDestinationModel

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

		resp.Diagnostics.AddError("Failed to read destination", detailFromError(getDestinationErr))

		return
	}

	model, modelDiags := NewCustomAPIDestinationModelFromResponse(ctx, getDestinationResponse.Response.Data)
	resp.Diagnostics.Append(modelDiags...)

	connectionDetails := sanitizedCustomAPIConnectionDetails(model.ConnectionDetails.Value(), state.Credentials.Value().CustomHeaders, nil)
	credentials := state.Credentials.Value()
	credentials.UpdateWithConnectionDetails(connectionDetails)

	model.Credentials = NewTypedObject(credentials)
	model.ConnectionDetails = NewTypedObject(connectionDetails)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.Identity.SetAttribute(ctx, path.Root("id"), model.ID)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func (r *customAPIDestinationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state, plan, config CustomAPIDestinationModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	resp.Diagnostics.Append(hydrateCustomAPIHeaderWriteOnlyValues(ctx, req.Config, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	params := cm.UpdateDestinationParams{
		DestinationID: plan.ID.ValueString(),
	}

	requestModel, writeOnlyValues, requestModelDiags := customAPIDestinationModelWithWriteOnlyCredentials(plan, config)
	resp.Diagnostics.Append(requestModelDiags...)

	if resp.Diagnostics.HasError() {
		return
	}

	updateDestinationRequest, updateDestinationRequestDiags := requestModel.ToUpdateDestinationData(ctx)
	resp.Diagnostics.Append(updateDestinationRequestDiags...)

	if resp.Diagnostics.HasError() {
		return
	}

	updateDestinationResponse, updateDestinationErr := r.providerData.client.UpdateDestination(ctx, &updateDestinationRequest, params)

	tflog.Info(ctx, "destination.update", map[string]any{
		"params":   params,
		"response": updateDestinationResponse,
		"err":      updateDestinationErr,
	})

	if updateDestinationResponse == nil {
		resp.Diagnostics.AddError("Failed to update destination", detailFromError(updateDestinationErr))

		return
	}

	model, modelDiags := NewCustomAPIDestinationModelFromResponse(ctx, updateDestinationResponse.Response.Data)
	resp.Diagnostics.Append(modelDiags...)

	model = sanitizedCustomAPIDestinationCredentials(model, plan, config, model.ConnectionDetails.Value(), writeOnlyValues)

	resp.Diagnostics.Append(resp.Identity.SetAttribute(ctx, path.Root("id"), model.ID)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
	resp.Diagnostics.Append(writeWriteOnlyCredentialVerifiers(ctx, resp.Private, writeOnlyValues)...)
}

//nolint:dupl
func (r *customAPIDestinationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state CustomAPIDestinationModel

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
