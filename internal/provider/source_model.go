package provider

import (
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type SourceModel struct {
	ID          types.Int64          `tfsdk:"id"`
	Type        types.String         `tfsdk:"type"`
	Label       types.String         `tfsdk:"label"`
	Credentials jsontypes.Normalized `tfsdk:"credentials"`
}
