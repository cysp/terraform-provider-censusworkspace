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
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = (*syncResource)(nil)
	_ resource.ResourceWithConfigure   = (*syncResource)(nil)
	_ resource.ResourceWithIdentity    = (*syncResource)(nil)
	_ resource.ResourceWithImportState = (*syncResource)(nil)
)

//nolint:ireturn
func NewSyncResource() resource.Resource {
	return &syncResource{}
}

type syncResource struct {
	providerData ProviderData
}

func (r *syncResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sync"
}

func (r *syncResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = SyncResourceSchema(ctx)
}

func (r *syncResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	resp.Diagnostics.Append(SetProviderDataFromResourceConfigureRequest(req, &r.providerData)...)
}

func (r *syncResource) IdentitySchema(ctx context.Context, _ resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = SyncResourceIdentitySchema(ctx)
}

func (r *syncResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	ImportStatePassthroughMultipartID(ctx, []path.Path{path.Root("id")}, req, resp)
}

func (r *syncResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan SyncModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	createSyncRequest, createSyncRequestDiags := plan.ToCreateOrUpdateSyncData(ctx)
	resp.Diagnostics.Append(createSyncRequestDiags...)

	if resp.Diagnostics.HasError() {
		return
	}

	createSyncResponse, createSyncErr := r.providerData.client.CreateSync(ctx, &createSyncRequest)

	tflog.Info(ctx, "sync.create", map[string]any{
		"request":  createSyncRequest,
		"response": createSyncResponse,
		"err":      createSyncErr,
	})

	if createSyncResponse == nil {
		resp.Diagnostics.AddError("Failed to create sync", createSyncErr.Error())

		return
	}

	syncID := createSyncResponse.Response.Data.SyncID
	syncIDString := strconv.FormatInt(syncID, 10)

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), syncIDString)...)

	if resp.Diagnostics.HasError() {
		return
	}

	getSyncParams := cm.GetSyncParams{
		SyncID: syncIDString,
	}

	getSyncResponse, getSyncErr := r.providerData.client.GetSync(ctx, getSyncParams)

	tflog.Info(ctx, "sync.read", map[string]any{
		"request":  getSyncParams,
		"response": getSyncResponse,
		"err":      getSyncErr,
	})

	model, modelDiags := NewSyncModelFromResponse(ctx, getSyncResponse.Response.Data)
	resp.Diagnostics.Append(modelDiags...)

	// if model.SyncEngine.IsNull() && !plan.SyncEngine.IsUnknown() {
	// 	model.SyncEngine = plan.SyncEngine
	// }

	// model.Credentials = plan.Credentials

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.Identity.SetAttribute(ctx, path.Root("id"), model.ID)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func (r *syncResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state SyncModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	params := cm.GetSyncParams{
		SyncID: state.ID.ValueString(),
	}

	getSyncResponse, getSyncErr := r.providerData.client.GetSync(ctx, params)

	tflog.Info(ctx, "sync.read", map[string]any{
		"params":   params,
		"response": getSyncResponse,
		"err":      getSyncErr,
	})

	if getSyncResponse == nil {
		var srsc *cm.StatusResponseStatusCode
		if errors.As(getSyncErr, &srsc) {
			if srsc.StatusCode == http.StatusNotFound {
				resp.Diagnostics.AddWarning("Failed to read sync", srsc.Error())
				resp.State.RemoveResource(ctx)

				return
			}
		}

		resp.Diagnostics.AddError("Failed to read sync", getSyncErr.Error())

		return
	}

	model, modelDiags := NewSyncModelFromResponse(ctx, getSyncResponse.Response.Data)
	resp.Diagnostics.Append(modelDiags...)

	// if model.SyncEngine.IsNull() && !state.SyncEngine.IsUnknown() {
	// 	model.SyncEngine = state.SyncEngine
	// }

	// model.Credentials = state.Credentials

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.Identity.SetAttribute(ctx, path.Root("id"), model.ID)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func (r *syncResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state, plan SyncModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	params := cm.UpdateSyncParams{
		SyncID: plan.ID.ValueString(),
	}

	updateSyncRequest, updateSyncRequestDiags := plan.ToCreateOrUpdateSyncData(ctx)
	resp.Diagnostics.Append(updateSyncRequestDiags...)

	if resp.Diagnostics.HasError() {
		return
	}

	updateSyncResponse, updateSyncErr := r.providerData.client.UpdateSync(ctx, &updateSyncRequest, params)

	tflog.Info(ctx, "sync.update", map[string]any{
		"params":   params,
		"request":  updateSyncRequest,
		"response": updateSyncResponse,
		"err":      updateSyncErr,
	})

	if updateSyncResponse == nil {
		resp.Diagnostics.AddError("Failed to update sync", updateSyncErr.Error())

		return
	}

	model, modelDiags := NewSyncModelFromResponse(ctx, updateSyncResponse.Response.Data)
	resp.Diagnostics.Append(modelDiags...)

	// if model.SyncEngine.IsNull() && !state.SyncEngine.IsUnknown() {
	// 	model.SyncEngine = state.SyncEngine
	// }

	// model.Credentials = plan.Credentials

	resp.Diagnostics.Append(resp.Identity.SetAttribute(ctx, path.Root("id"), model.ID)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

//nolint:dupl
func (r *syncResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state SyncModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	params := cm.DeleteSyncParams{
		SyncID: state.ID.ValueString(),
	}

	deleteSyncResponse, deleteSyncErr := r.providerData.client.DeleteSync(ctx, params)

	tflog.Info(ctx, "sync.delete", map[string]any{
		"params":   params,
		"response": deleteSyncResponse,
		"err":      deleteSyncErr,
	})

	var srsc *cm.StatusResponseStatusCode
	if errors.As(deleteSyncErr, &srsc) {
		if srsc.StatusCode == http.StatusNotFound {
			resp.Diagnostics.AddWarning("Sync not found", srsc.Error())
			resp.State.RemoveResource(ctx)

			return
		}
	}

	if deleteSyncResponse == nil || deleteSyncResponse.Response.Status.ResponseStatus != cm.ResponseStatusDeleted {
		var detail string

		if deleteSyncResponse != nil {
			detail = deleteSyncResponse.Response.Message.Value
		}

		if detail == "" && deleteSyncErr != nil {
			detail = deleteSyncErr.Error()
		}

		resp.Diagnostics.AddError("Failed to delete sync", detail)

		return
	}
}
