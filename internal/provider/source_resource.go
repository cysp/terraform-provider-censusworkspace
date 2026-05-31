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
	ImportStatePassthroughMultipartID(ctx, []path.Path{path.Root("id")}, req, resp)
}

func (r *sourceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	r.managedResource().Create(ctx, req, resp)
}

func (r *sourceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	r.managedResource().Read(ctx, req, resp)
}

func (r *sourceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	r.managedResource().Update(ctx, req, resp)
}

func (r *sourceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	r.managedResource().Delete(ctx, req, resp)
}

func (r *sourceResource) managedResource() managedResource[SourceModel, cm.CreateSourceBody, cm.UpdateSourceBody, cm.SourceData] {
	return managedResource[SourceModel, cm.CreateSourceBody, cm.UpdateSourceBody, cm.SourceData]{
		readErrorTitle:    "Failed to read source",
		createErrorTitle:  "Failed to create source",
		updateErrorTitle:  "Failed to update source",
		deleteErrorTitle:  "Failed to delete source",
		deleteMissingText: "Source not found",
		createRequest: func(ctx context.Context, model SourceModel) (cm.CreateSourceBody, diag.Diagnostics) {
			return model.ToCreateSourceData(ctx)
		},
		updateRequest: func(ctx context.Context, model SourceModel) (cm.UpdateSourceBody, diag.Diagnostics) {
			return model.ToUpdateSourceData(ctx)
		},
		modelFromRead: NewSourceModelFromResponse,
		afterCreateRead: func(plan SourceModel, model *SourceModel) {
			if model.SyncEngine.IsNull() && !plan.SyncEngine.IsUnknown() {
				model.SyncEngine = plan.SyncEngine
			}

			model.Credentials = plan.Credentials
		},
		afterRead: func(state SourceModel, model *SourceModel) {
			if model.SyncEngine.IsNull() && !state.SyncEngine.IsUnknown() {
				model.SyncEngine = state.SyncEngine
			}

			model.Credentials = state.Credentials
		},
		afterUpdate: func(state SourceModel, plan SourceModel, model *SourceModel) {
			if model.SyncEngine.IsNull() && !state.SyncEngine.IsUnknown() {
				model.SyncEngine = state.SyncEngine
			}

			model.Credentials = plan.Credentials
		},
		create: createResourceOperation[cm.CreateSourceBody]{
			name: "source.create",
			run: func(ctx context.Context, request cm.CreateSourceBody) (int64, error) {
				response, err := r.providerData.client.CreateSource(ctx, &request)
				if err != nil {
					return 0, fmt.Errorf("create source: %w", err)
				}

				if response == nil {
					return 0, responseMissing("source.create")
				}

				return response.Response.Data.ID, nil
			},
		},
		read: resourceOperation[string, cm.SourceData]{
			name: "source.read",
			run: func(ctx context.Context, id string) (cm.SourceData, error) {
				response, err := r.providerData.client.GetSource(ctx, cm.GetSourceParams{SourceID: id})
				if err != nil {
					return cm.SourceData{}, fmt.Errorf("read source: %w", err)
				}

				if response == nil {
					return cm.SourceData{}, responseMissing("source.read")
				}

				return response.Response.Data, nil
			},
		},
		update: resourceOperation[updateResourceRequest[cm.UpdateSourceBody], cm.SourceData]{
			name: "source.update",
			run: func(ctx context.Context, request updateResourceRequest[cm.UpdateSourceBody]) (cm.SourceData, error) {
				response, err := r.providerData.client.UpdateSource(ctx, &request.request, cm.UpdateSourceParams{SourceID: request.id})
				if err != nil {
					return cm.SourceData{}, fmt.Errorf("update source: %w", err)
				}

				if response == nil {
					return cm.SourceData{}, responseMissing("source.update")
				}

				return response.Response.Data, nil
			},
		},
		delete: deleteResourceOperation{
			name: "source.delete",
			run: func(ctx context.Context, id string) (*cm.StatusResponseStatusCode, error) {
				return r.providerData.client.DeleteSource(ctx, cm.DeleteSourceParams{SourceID: id})
			},
		},
	}
}
