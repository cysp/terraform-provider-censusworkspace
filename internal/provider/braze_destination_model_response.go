//nolint:dupl
package provider

import (
	"context"
	"strconv"

	cm "github.com/cysp/terraform-provider-censusworkspace/internal/census-management-go"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewBrazeDestinationModelFromResponse(ctx context.Context, response cm.DestinationData) (BrazeDestinationModel, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	path := path.Empty()

	model := BrazeDestinationModel{
		destinationModelBase: destinationModelBase{
			ID:        types.StringValue(strconv.FormatInt(response.ID, 10)),
			Name:      types.StringValue(response.Name),
			CreatedAt: timetypes.NewRFC3339TimeValue(response.CreatedAt),
		},
	}

	if response.ConnectionDetails != nil {
		path := path.AtName("connection_details")

		connectionDetails, connectionDetailsDiags := NewBrazeDestinationConnectionDetailsFromResponse(ctx, path, response.ConnectionDetails)
		diags.Append(connectionDetailsDiags...)

		model.ConnectionDetails = connectionDetails
	}

	if lastTestedAt, lastTestedAtOk := response.LastTestedAt.Get(); lastTestedAtOk {
		model.LastTestedAt = timetypes.NewRFC3339TimeValue(lastTestedAt)
	}

	if lastTestSucceeded, lastTestSucceededOk := response.LastTestSucceeded.Get(); lastTestSucceededOk {
		model.LastTestSucceeded = types.BoolValue(lastTestSucceeded)
	}

	return model, diags
}
