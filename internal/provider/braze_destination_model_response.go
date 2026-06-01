package provider

import (
	"context"
	"strconv"

	cm "github.com/cysp/terraform-provider-censusworkspace/internal/census-management-go"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewBrazeDestinationModelFromResponse(ctx context.Context, response cm.DestinationData) (BrazeDestinationModel, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	path := path.Empty()

	model := BrazeDestinationModel{
		destinationModelBase: destinationModelBase{
			ID:   types.StringValue(strconv.FormatInt(response.ID, 10)),
			Name: types.StringValue(response.Name),
		},
	}

	if response.ConnectionDetails != nil {
		path := path.AtName("connection_details")

		connectionDetails, connectionDetailsDiags := NewBrazeDestinationConnectionDetailsFromResponse(ctx, path, response.ConnectionDetails)
		diags.Append(connectionDetailsDiags...)

		model.ConnectionDetails = connectionDetails
	}

	return model, diags
}
