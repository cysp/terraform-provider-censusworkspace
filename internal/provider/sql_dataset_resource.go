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
	var plan SQLDatasetModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	createDatasetRequest, createDatasetRequestDiags := plan.ToCreateDatasetBody(ctx)
	resp.Diagnostics.Append(createDatasetRequestDiags...)

	if resp.Diagnostics.HasError() {
		return
	}

	createDatasetResponse, createDatasetErr := r.providerData.client.CreateDataset(ctx, createDatasetRequest)

	tflog.Info(ctx, "sql_dataset.create", map[string]any{
		"request":  createDatasetRequest,
		"response": createDatasetResponse,
		"err":      createDatasetErr,
	})

	if createDatasetResponse == nil {
		resp.Diagnostics.AddError("Failed to create dataset", createDatasetErr.Error())

		return
	}

	datasetID := createDatasetResponse.Response.Data.ID
	datasetIDString := strconv.FormatInt(datasetID, 10)

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), datasetIDString)...)

	if resp.Diagnostics.HasError() {
		return
	}

	getDatasetParams := cm.GetDatasetParams{
		DatasetID: datasetIDString,
	}

	getDatasetResponse, getDatasetErr := r.providerData.client.GetDataset(ctx, getDatasetParams)

	tflog.Info(ctx, "sql_dataset.read", map[string]any{
		"request":  getDatasetParams,
		"response": getDatasetResponse,
		"err":      getDatasetErr,
	})

	model, modelDiags := NewSQLDatasetModelFromResponse(ctx, path.Empty(), getDatasetResponse.Response.Data)
	resp.Diagnostics.Append(modelDiags...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.Identity.SetAttribute(ctx, path.Root("id"), model.ID)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func (r *sqlDatasetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state SQLDatasetModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	params := cm.GetDatasetParams{
		DatasetID: state.ID.ValueString(),
	}

	getDatasetResponse, getDatasetErr := r.providerData.client.GetDataset(ctx, params)

	tflog.Info(ctx, "sql_dataset.read", map[string]any{
		"params":   params,
		"response": getDatasetResponse,
		"err":      getDatasetErr,
	})

	if getDatasetResponse == nil {
		var srsc *cm.StatusResponseStatusCode
		if errors.As(getDatasetErr, &srsc) {
			if srsc.StatusCode == http.StatusNotFound {
				resp.Diagnostics.AddWarning("Failed to read dataset", srsc.Error())
				resp.State.RemoveResource(ctx)

				return
			}
		}

		resp.Diagnostics.AddError("Failed to read dataset", getDatasetErr.Error())

		return
	}

	model, modelDiags := NewSQLDatasetModelFromResponse(ctx, path.Empty(), getDatasetResponse.Response.Data)
	resp.Diagnostics.Append(modelDiags...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.Identity.SetAttribute(ctx, path.Root("id"), model.ID)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func (r *sqlDatasetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state, plan SQLDatasetModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	params := cm.UpdateDatasetParams{
		DatasetID: plan.ID.ValueString(),
	}

	updateDatasetRequest, updateDatasetRequestDiags := plan.ToUpdateDatasetBody(ctx)
	resp.Diagnostics.Append(updateDatasetRequestDiags...)

	if resp.Diagnostics.HasError() {
		return
	}

	updateDatasetResponse, updateDatasetErr := r.providerData.client.UpdateDataset(ctx, updateDatasetRequest, params)

	tflog.Info(ctx, "sql_dataset.update", map[string]any{
		"params":   params,
		"request":  updateDatasetRequest,
		"response": updateDatasetResponse,
		"err":      updateDatasetErr,
	})

	if updateDatasetResponse == nil {
		resp.Diagnostics.AddError("Failed to update dataset", updateDatasetErr.Error())

		return
	}

	model, modelDiags := NewSQLDatasetModelFromResponse(ctx, path.Empty(), updateDatasetResponse.Response.Data)
	resp.Diagnostics.Append(modelDiags...)

	resp.Diagnostics.Append(resp.Identity.SetAttribute(ctx, path.Root("id"), model.ID)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

//nolint:dupl
func (r *sqlDatasetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state SQLDatasetModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	params := cm.DeleteDatasetParams{
		DatasetID: state.ID.ValueString(),
	}

	deleteDatasetResponse, deleteDatasetErr := r.providerData.client.DeleteDataset(ctx, params)

	tflog.Info(ctx, "sql_dataset.delete", map[string]any{
		"params":   params,
		"response": deleteDatasetResponse,
		"err":      deleteDatasetErr,
	})

	var srsc *cm.StatusResponseStatusCode
	if errors.As(deleteDatasetErr, &srsc) {
		if srsc.StatusCode == http.StatusNotFound {
			resp.Diagnostics.AddWarning("dataset not found", srsc.Error())
			resp.State.RemoveResource(ctx)

			return
		}
	}

	if deleteDatasetResponse == nil || deleteDatasetResponse.Response.Status.ResponseStatus != cm.ResponseStatusDeleted {
		var detail string

		if deleteDatasetResponse != nil {
			detail = deleteDatasetResponse.Response.Message.Value
		}

		if detail == "" && deleteDatasetErr != nil {
			detail = deleteDatasetErr.Error()
		}

		resp.Diagnostics.AddError("Failed to delete dataset", detail)

		return
	}
}
