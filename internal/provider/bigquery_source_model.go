package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	BigQuerySourceType = "big_query"
)

type BigQuerySourceModel struct {
	sourceModelCommon
	Credentials       BigQuerySourceCredentials       `tfsdk:"credentials"`
	ConnectionDetails BigQuerySourceConnectionDetails `tfsdk:"connection_details"`
}

type BigQuerySourceCredentials struct {
	ProjectID types.String `tfsdk:"project_id"`
	Location  types.String `tfsdk:"location"`
}
