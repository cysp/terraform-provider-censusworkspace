package provider

import (
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type SourceModel struct {
	sourceModelBase

	Type              types.String         `tfsdk:"type"`
	Credentials       jsontypes.Normalized `tfsdk:"credentials"`
	ConnectionDetails jsontypes.Normalized `tfsdk:"connection_details"`
}
