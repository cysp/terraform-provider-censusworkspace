package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	BigQuerySourceType = "big_query"
)

type BigQuerySourceModel struct {
	sourceModelBase

	Credentials       TypedObject[BigQuerySourceCredentials]       `tfsdk:"credentials"`
	ConnectionDetails TypedObject[BigQuerySourceConnectionDetails] `tfsdk:"connection_details"`
}

type BigQuerySourceCredentials struct {
	ProjectID         types.String `tfsdk:"project_id"`
	Location          types.String `tfsdk:"location"`
	ServiceAccountKey types.String `tfsdk:"service_account_key"`
}

type BigQuerySourceConnectionDetails struct {
	ProjectID      types.String `tfsdk:"project_id"`
	Location       types.String `tfsdk:"location"`
	ServiceAccount types.String `tfsdk:"service_account"`
}
