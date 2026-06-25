package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	CustomAPIDestinationType = "custom_api"
)

type CustomAPIDestinationModel struct {
	destinationModelBase

	Credentials       TypedObject[CustomAPIDestinationCredentials]       `tfsdk:"credentials"`
	ConnectionDetails TypedObject[CustomAPIDestinationConnectionDetails] `tfsdk:"connection_details"`
}

//nolint:recvcheck
type CustomAPIDestinationCredentials struct {
	APIVersion    types.Int64                                             `tfsdk:"api_version"`
	WebhookURL    types.String                                            `tfsdk:"webhook_url"`
	CustomHeaders TypedMap[TypedObject[CustomAPIDestinationCustomHeader]] `tfsdk:"custom_headers"`
}

type CustomAPIDestinationConnectionDetails struct {
	APIVersion    types.Int64                                                              `tfsdk:"api_version"`
	WebhookURL    types.String                                                             `tfsdk:"webhook_url"`
	CustomHeaders TypedMap[TypedObject[CustomAPIDestinationConnectionDetailsCustomHeader]] `tfsdk:"custom_headers"`
}

type CustomAPIDestinationCustomHeader struct {
	Value    types.String `tfsdk:"value"`
	ValueWO  types.String `tfsdk:"value_wo"`
	IsSecret types.Bool   `tfsdk:"is_secret"`
}

type CustomAPIDestinationConnectionDetailsCustomHeader struct {
	Value    types.String `tfsdk:"value"`
	IsSecret types.Bool   `tfsdk:"is_secret"`
}
