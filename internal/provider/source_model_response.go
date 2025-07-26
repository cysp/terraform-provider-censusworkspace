package provider

import (
	"context"

	cm "github.com/cysp/terraform-provider-censusworkspace/internal/census-management-go/client"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewSourceResourceModelFromResponse(_ context.Context, response cm.SourceData) (SourceModel, diag.Diagnostics) {
	model := SourceModel{
		ID:    types.Int64Value(response.ID),
		Type:  types.StringValue(response.Type),
		Label: types.StringPointerValue(response.Label.ValueStringPointer()),
	}

	return model, nil
}
