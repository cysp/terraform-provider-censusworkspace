package provider

import (
	"context"
	"strconv"

	cm "github.com/cysp/terraform-provider-censusworkspace/internal/census-management-go"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewSyncModelFromResponse(ctx context.Context, sync cm.SyncData) (SyncModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	model := SyncModel{
		ID:        types.StringValue(strconv.FormatInt(sync.ID, 10)),
		Label:     types.StringPointerValue(sync.Label.ValueStringPointer()),
		Operation: types.StringValue(sync.Operation),
	}

	return model, diags
}
