package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	BigQueryDestinationType = "big_query"
)

type BigQueryDestinationModel struct {
	destinationModelBase

	Credentials       TypedObject[BigQueryDestinationCredentials]       `tfsdk:"credentials"`
	ConnectionDetails TypedObject[BigQueryDestinationConnectionDetails] `tfsdk:"connection_details"`
}

//nolint:recvcheck
type BigQueryDestinationCredentials struct {
	ProjectID         types.String `tfsdk:"project_id"`
	Location          types.String `tfsdk:"location"`
	ServiceAccountKey types.String `tfsdk:"service_account_key"`
}

type BigQueryDestinationConnectionDetails struct {
	ProjectID           types.String `tfsdk:"project_id"`
	Location            types.String `tfsdk:"location"`
	ServiceAccountEmail types.String `tfsdk:"service_account_email"`
	ServiceAccountKey   types.String `tfsdk:"service_account_key"`
}
