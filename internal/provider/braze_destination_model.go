package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	BrazeDestinationType = "braze"
)

type BrazeDestinationModel struct {
	destinationModelBase

	Credentials       TypedObject[BrazeDestinationCredentials]       `tfsdk:"credentials"`
	ConnectionDetails TypedObject[BrazeDestinationConnectionDetails] `tfsdk:"connection_details"`
}

//nolint:recvcheck
type BrazeDestinationCredentials struct {
	InstanceURL types.String `tfsdk:"instance_url"`
	APIKey      types.String `tfsdk:"api_key"`
	APIKeyWO    types.String `tfsdk:"api_key_wo"`
	ClientKey   types.String `tfsdk:"client_key"`
	ClientKeyWO types.String `tfsdk:"client_key_wo"`
}

type BrazeDestinationConnectionDetails struct {
	InstanceURL types.String `tfsdk:"instance_url"`
}
