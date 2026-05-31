package provider

import (
	"context"
	"fmt"

	cm "github.com/cysp/terraform-provider-censusworkspace/internal/census-management-go"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var (
	_ resource.Resource                = (*destinationResource)(nil)
	_ resource.ResourceWithConfigure   = (*destinationResource)(nil)
	_ resource.ResourceWithIdentity    = (*destinationResource)(nil)
	_ resource.ResourceWithImportState = (*destinationResource)(nil)
)

const (
	DestinationResourceTypeName = "censusworkspace_destination"
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
	ImportStatePassthroughMultipartID(ctx, []path.Path{path.Root("id")}, req, resp)
}

func (r *destinationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	r.managedResource().Create(ctx, req, resp)
}

func (r *destinationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	r.managedResource().Read(ctx, req, resp)
}

func (r *destinationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	r.managedResource().Update(ctx, req, resp)
}

func (r *destinationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	r.managedResource().Delete(ctx, req, resp)
}

func (r *destinationResource) managedResource() managedResource[DestinationModel, cm.CreateDestinationBody, cm.UpdateDestinationBody, cm.DestinationData] {
	return managedResource[DestinationModel, cm.CreateDestinationBody, cm.UpdateDestinationBody, cm.DestinationData]{
		readErrorTitle:    "Failed to read destination",
		createErrorTitle:  "Failed to create destination",
		updateErrorTitle:  "Failed to update destination",
		deleteErrorTitle:  "Failed to delete destination",
		deleteMissingText: "Destination not found",
		createRequest: func(ctx context.Context, model DestinationModel) (cm.CreateDestinationBody, diag.Diagnostics) {
			return model.ToCreateDestinationData(ctx)
		},
		updateRequest: func(ctx context.Context, model DestinationModel) (cm.UpdateDestinationBody, diag.Diagnostics) {
			return model.ToUpdateDestinationData(ctx)
		},
		modelFromRead: NewDestinationModelFromResponse,
		afterCreateRead: func(plan DestinationModel, model *DestinationModel) {
			model.Credentials = plan.Credentials
		},
		afterRead: func(state DestinationModel, model *DestinationModel) {
			model.Credentials = state.Credentials
		},
		afterUpdate: func(_ DestinationModel, plan DestinationModel, model *DestinationModel) {
			model.Credentials = plan.Credentials
		},
		create: createResourceOperation[cm.CreateDestinationBody]{
			name: "destination.create",
			run: func(ctx context.Context, request cm.CreateDestinationBody) (int64, error) {
				response, err := r.providerData.client.CreateDestination(ctx, &request)
				if err != nil {
					return 0, fmt.Errorf("create destination: %w", err)
				}

				if response == nil {
					return 0, responseMissing("destination.create")
				}

				return response.Response.Data.ID, nil
			},
		},
		read: resourceOperation[string, cm.DestinationData]{
			name: "destination.read",
			run: func(ctx context.Context, id string) (cm.DestinationData, error) {
				response, err := r.providerData.client.GetDestination(ctx, cm.GetDestinationParams{DestinationID: id})
				if err != nil {
					return cm.DestinationData{}, fmt.Errorf("read destination: %w", err)
				}

				if response == nil {
					return cm.DestinationData{}, responseMissing("destination.read")
				}

				return response.Response.Data, nil
			},
		},
		update: resourceOperation[updateResourceRequest[cm.UpdateDestinationBody], cm.DestinationData]{
			name: "destination.update",
			run: func(ctx context.Context, request updateResourceRequest[cm.UpdateDestinationBody]) (cm.DestinationData, error) {
				response, err := r.providerData.client.UpdateDestination(ctx, &request.request, cm.UpdateDestinationParams{DestinationID: request.id})
				if err != nil {
					return cm.DestinationData{}, fmt.Errorf("update destination: %w", err)
				}

				if response == nil {
					return cm.DestinationData{}, responseMissing("destination.update")
				}

				return response.Response.Data, nil
			},
		},
		delete: deleteResourceOperation{
			name: "destination.delete",
			run: func(ctx context.Context, id string) (*cm.StatusResponseStatusCode, error) {
				return r.providerData.client.DeleteDestination(ctx, cm.DeleteDestinationParams{DestinationID: id})
			},
		},
	}
}
