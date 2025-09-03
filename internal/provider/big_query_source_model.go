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

//nolint:recvcheck
type BigQuerySourceCredentials struct {
	ProjectID         types.String                                            `tfsdk:"project_id"`
	Location          types.String                                            `tfsdk:"location"`
	ServiceAccountKey TypedObject[BigQuerySourceCredentialsServiceAccountKey] `tfsdk:"service_account_key"`
}

type BigQuerySourceCredentialsServiceAccountKey struct {
	ProjectID    types.String `tfsdk:"project_id"`
	PrivateKeyID types.String `tfsdk:"private_key_id"`
	PrivateKey   types.String `tfsdk:"private_key"`
	ClientEmail  types.String `tfsdk:"client_email"`
	ClientID     types.String `tfsdk:"client_id"`
}

type BigQuerySourceConnectionDetails struct {
	ProjectID      types.String `tfsdk:"project_id"`
	Location       types.String `tfsdk:"location"`
	ServiceAccount types.String `tfsdk:"service_account"`
}
