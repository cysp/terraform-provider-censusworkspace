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
	_ resource.Resource                = (*sqlDatasetResource)(nil)
	_ resource.ResourceWithConfigure   = (*sqlDatasetResource)(nil)
	_ resource.ResourceWithIdentity    = (*sqlDatasetResource)(nil)
	_ resource.ResourceWithImportState = (*sqlDatasetResource)(nil)
)

//nolint:ireturn
func NewSQLDatasetResource() resource.Resource {
	return &sqlDatasetResource{}
}

type sqlDatasetResource struct {
	providerData ProviderData
}

func (r *sqlDatasetResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sql_dataset"
}

func (r *sqlDatasetResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = SQLDatasetResourceSchema(ctx)
}

func (r *sqlDatasetResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	resp.Diagnostics.Append(SetProviderDataFromResourceConfigureRequest(req, &r.providerData)...)
}

func (r *sqlDatasetResource) IdentitySchema(ctx context.Context, _ resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = SQLDatasetResourceIdentitySchema(ctx)
}

func (r *sqlDatasetResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	ImportStatePassthroughMultipartID(ctx, []path.Path{path.Root("id")}, req, resp)
}

func (r *sqlDatasetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	r.managedResource().Create(ctx, req, resp)
}

func (r *sqlDatasetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	r.managedResource().Read(ctx, req, resp)
}

func (r *sqlDatasetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	r.managedResource().Update(ctx, req, resp)
}

func (r *sqlDatasetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	r.managedResource().Delete(ctx, req, resp)
}

func (r *sqlDatasetResource) managedResource() managedResource[SQLDatasetModel, cm.CreateDatasetBody, cm.UpdateDatasetBody, cm.DatasetData] {
	return managedResource[SQLDatasetModel, cm.CreateDatasetBody, cm.UpdateDatasetBody, cm.DatasetData]{
		readErrorTitle:    "Failed to read dataset",
		createErrorTitle:  "Failed to create dataset",
		updateErrorTitle:  "Failed to update dataset",
		deleteErrorTitle:  "Failed to delete dataset",
		deleteMissingText: "dataset not found",
		createRequest: func(ctx context.Context, model SQLDatasetModel) (cm.CreateDatasetBody, diag.Diagnostics) {
			return model.ToCreateDatasetBody(ctx)
		},
		updateRequest: func(ctx context.Context, model SQLDatasetModel) (cm.UpdateDatasetBody, diag.Diagnostics) {
			return model.ToUpdateDatasetBody(ctx)
		},
		modelFromRead: func(ctx context.Context, data cm.DatasetData) (SQLDatasetModel, diag.Diagnostics) {
			return NewSQLDatasetModelFromResponse(ctx, path.Empty(), data)
		},
		create: createResourceOperation[cm.CreateDatasetBody]{
			name: "sql_dataset.create",
			run: func(ctx context.Context, request cm.CreateDatasetBody) (int64, error) {
				response, err := r.providerData.client.CreateDataset(ctx, request)
				if err != nil {
					return 0, fmt.Errorf("create dataset: %w", err)
				}

				if response == nil {
					return 0, responseMissing("sql_dataset.create")
				}

				return response.Response.Data.ID, nil
			},
		},
		read: resourceOperation[string, cm.DatasetData]{
			name: "sql_dataset.read",
			run: func(ctx context.Context, id string) (cm.DatasetData, error) {
				response, err := r.providerData.client.GetDataset(ctx, cm.GetDatasetParams{DatasetID: id})
				if err != nil {
					return cm.DatasetData{}, fmt.Errorf("read dataset: %w", err)
				}

				if response == nil {
					return cm.DatasetData{}, responseMissing("sql_dataset.read")
				}

				return response.Response.Data, nil
			},
		},
		update: resourceOperation[updateResourceRequest[cm.UpdateDatasetBody], cm.DatasetData]{
			name: "sql_dataset.update",
			run: func(ctx context.Context, request updateResourceRequest[cm.UpdateDatasetBody]) (cm.DatasetData, error) {
				response, err := r.providerData.client.UpdateDataset(ctx, request.request, cm.UpdateDatasetParams{DatasetID: request.id})
				if err != nil {
					return cm.DatasetData{}, fmt.Errorf("update dataset: %w", err)
				}

				if response == nil {
					return cm.DatasetData{}, responseMissing("sql_dataset.update")
				}

				return response.Response.Data, nil
			},
		},
		delete: deleteResourceOperation{
			name: "sql_dataset.delete",
			run: func(ctx context.Context, id string) (*cm.StatusResponseStatusCode, error) {
				return r.providerData.client.DeleteDataset(ctx, cm.DeleteDatasetParams{DatasetID: id})
			},
		},
	}
}
