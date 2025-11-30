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

func NewBigQueryDestinationModelFromResponse(ctx context.Context, response cm.DestinationData) (BigQueryDestinationModel, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	model := BigQueryDestinationModel{
		destinationModelBase: destinationModelBase{
			ID:        types.StringValue(strconv.FormatInt(response.ID, 10)),
			Name:      types.StringValue(response.Name),
			CreatedAt: timetypes.NewRFC3339TimeValue(response.CreatedAt),
		},
	}

	if response.ConnectionDetails != nil {
		connectionDetailsPath := path.Root("connection_details")

		connectionDetails, connectionDetailsDiags := NewBigQueryDestinationConnectionDetailsFromResponse(ctx, connectionDetailsPath, response.ConnectionDetails)
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
