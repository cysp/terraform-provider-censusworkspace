package provider

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	cm "github.com/cysp/terraform-provider-censusworkspace/internal/census-management-go"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var errNilOperationResponse = errors.New("operation returned nil response")

type resourceModel interface {
	getID() types.String
}

type resourceOperation[Request any, Response any] struct {
	name string
	run  func(context.Context, Request) (Response, error)
}

type createResourceOperation[Request any] struct {
	name string
	run  func(context.Context, Request) (int64, error)
}

type deleteResourceOperation struct {
	name string
	run  func(context.Context, string) (*cm.StatusResponseStatusCode, error)
}

type managedResource[
	Model resourceModel,
	CreateRequest any,
	UpdateRequest any,
	ReadResponse any,
] struct {
	readErrorTitle    string
	createErrorTitle  string
	updateErrorTitle  string
	deleteErrorTitle  string
	deleteMissingText string

	createRequest func(context.Context, Model) (CreateRequest, diag.Diagnostics)
	updateRequest func(context.Context, Model) (UpdateRequest, diag.Diagnostics)
	modelFromRead func(context.Context, ReadResponse) (Model, diag.Diagnostics)

	afterCreateRead func(plan Model, model *Model)
	afterRead       func(state Model, model *Model)
	afterUpdate     func(state Model, plan Model, model *Model)

	create createResourceOperation[CreateRequest]
	read   resourceOperation[string, ReadResponse]
	update resourceOperation[updateResourceRequest[UpdateRequest], ReadResponse]
	delete deleteResourceOperation
}

type updateResourceRequest[Request any] struct {
	id      string
	request Request
}

func (m managedResource[Model, CreateRequest, UpdateRequest, ReadResponse]) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan Model

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	createRequest, createRequestDiags := m.createRequest(ctx, plan)
	resp.Diagnostics.Append(createRequestDiags...)

	if resp.Diagnostics.HasError() {
		return
	}

	id, ok := m.runCreate(ctx, createRequest, &resp.Diagnostics)
	if !ok {
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)

	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, readErr := m.runRead(ctx, id)
	if readErr != nil {
		resp.Diagnostics.AddError(m.createErrorTitle, detailFromError(readErr))

		return
	}

	model, modelDiags := m.modelFromRead(ctx, readResponse)
	resp.Diagnostics.Append(modelDiags...)

	if m.afterCreateRead != nil {
		m.afterCreateRead(plan, &model)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	m.setIdentityAndState(ctx, model, &resp.Diagnostics, resp.Identity.SetAttribute, resp.State.Set)
}

func (m managedResource[Model, CreateRequest, UpdateRequest, ReadResponse]) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state Model

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, readErr := m.runRead(ctx, state.getID().ValueString())
	if readErr != nil {
		if isNotFound(readErr) {
			resp.Diagnostics.AddWarning(m.readErrorTitle, readErr.Error())
			resp.State.RemoveResource(ctx)

			return
		}

		resp.Diagnostics.AddError(m.readErrorTitle, detailFromError(readErr))

		return
	}

	model, modelDiags := m.modelFromRead(ctx, readResponse)
	resp.Diagnostics.Append(modelDiags...)

	if m.afterRead != nil {
		m.afterRead(state, &model)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	m.setIdentityAndState(ctx, model, &resp.Diagnostics, resp.Identity.SetAttribute, resp.State.Set)
}

func (m managedResource[Model, CreateRequest, UpdateRequest, ReadResponse]) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state, plan Model

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	updateRequest, updateRequestDiags := m.updateRequest(ctx, plan)
	resp.Diagnostics.Append(updateRequestDiags...)

	if resp.Diagnostics.HasError() {
		return
	}

	updateResponse, updateErr := m.update.run(ctx, updateResourceRequest[UpdateRequest]{
		id:      plan.getID().ValueString(),
		request: updateRequest,
	})
	tflog.Info(ctx, m.update.name, map[string]any{
		"id":       plan.getID().ValueString(),
		"request":  updateRequest,
		"response": updateResponse,
		"err":      updateErr,
	})

	if updateErr != nil {
		resp.Diagnostics.AddError(m.updateErrorTitle, detailFromError(updateErr))

		return
	}

	model, modelDiags := m.modelFromRead(ctx, updateResponse)
	resp.Diagnostics.Append(modelDiags...)

	if m.afterUpdate != nil {
		m.afterUpdate(state, plan, &model)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	m.setIdentityAndState(ctx, model, &resp.Diagnostics, resp.Identity.SetAttribute, resp.State.Set)
}

func (m managedResource[Model, CreateRequest, UpdateRequest, ReadResponse]) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state Model

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	id := state.getID().ValueString()
	deleteResponse, deleteErr := m.delete.run(ctx, id)
	tflog.Info(ctx, m.delete.name, map[string]any{
		"id":       id,
		"response": deleteResponse,
		"err":      deleteErr,
	})

	if isNotFound(deleteErr) {
		resp.Diagnostics.AddWarning(m.deleteMissingText, deleteErr.Error())
		resp.State.RemoveResource(ctx)

		return
	}

	if deleteResponse == nil {
		resp.Diagnostics.AddError(m.deleteErrorTitle, responseMissing(m.delete.name).Error())

		return
	}

	if deleteResponse.Response.Status.ResponseStatus != cm.ResponseStatusDeleted {
		resp.Diagnostics.AddError(m.deleteErrorTitle, deleteDetail(deleteResponse, deleteErr))

		return
	}
}

func (m managedResource[Model, CreateRequest, UpdateRequest, ReadResponse]) runCreate(ctx context.Context, request CreateRequest, diags *diag.Diagnostics) (string, bool) {
	id, err := m.create.run(ctx, request)
	tflog.Info(ctx, m.create.name, map[string]any{
		"request": request,
		"id":      id,
		"err":     err,
	})

	if err != nil {
		diags.AddError(m.createErrorTitle, detailFromError(err))

		return "", false
	}

	return strconv.FormatInt(id, 10), true
}

//nolint:ireturn
func (m managedResource[Model, CreateRequest, UpdateRequest, ReadResponse]) runRead(ctx context.Context, id string) (ReadResponse, error) {
	response, err := m.read.run(ctx, id)
	tflog.Info(ctx, m.read.name, map[string]any{
		"id":       id,
		"response": response,
		"err":      err,
	})

	if err != nil {
		var zero ReadResponse

		return zero, err
	}

	return response, nil
}

func (m managedResource[Model, CreateRequest, UpdateRequest, ReadResponse]) setIdentityAndState(
	ctx context.Context,
	model Model,
	diags *diag.Diagnostics,
	setIdentity func(context.Context, path.Path, any) diag.Diagnostics,
	setState func(context.Context, any) diag.Diagnostics,
) {
	diags.Append(setIdentity(ctx, path.Root("id"), model.getID())...)

	if diags.HasError() {
		return
	}

	diags.Append(setState(ctx, &model)...)
}

func isNotFound(err error) bool {
	var statusResponse *cm.StatusResponseStatusCode

	return errors.As(err, &statusResponse) && statusResponse.StatusCode == http.StatusNotFound
}

func deleteDetail(response *cm.StatusResponseStatusCode, err error) string {
	if response != nil {
		if detail := response.Response.Message.Value; detail != "" {
			return detail
		}
	}

	return detailFromError(err)
}

func responseMissing(operation string) error {
	return fmt.Errorf("%s: %w", operation, errNilOperationResponse)
}
