package provider

import (
	"context"
	"strconv"

	cm "github.com/cysp/terraform-provider-censusworkspace/internal/census-management-go"
	"github.com/go-faster/jx"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewCustomAPIDestinationModelFromResponse(_ context.Context, response cm.DestinationData) (CustomAPIDestinationModel, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	model := CustomAPIDestinationModel{
		destinationModelBase: destinationModelBase{
			ID:        types.StringValue(strconv.FormatInt(response.ID, 10)),
			Name:      types.StringValue(response.Name),
			CreatedAt: timetypes.NewRFC3339TimeValue(response.CreatedAt),
		},
	}

	if response.ConnectionDetails != nil {
		dec := jx.DecodeBytes(response.ConnectionDetails)

		connectionDetailsModel := CustomAPIDestinationConnectionDetails{}

		connectionDetailsDecodeErr := connectionDetailsModel.Decode(dec)
		if connectionDetailsDecodeErr != nil {
			diags.AddAttributeError(path.Root("connection_details"), "Failed to decode value", connectionDetailsDecodeErr.Error())
		}

		model.ConnectionDetails = NewTypedObject(connectionDetailsModel)
	}

	if lastTestedAt, lastTestedAtOk := response.LastTestedAt.Get(); lastTestedAtOk {
		model.LastTestedAt = timetypes.NewRFC3339TimeValue(lastTestedAt)
	}

	if lastTestSucceeded, lastTestSucceededOk := response.LastTestSucceeded.Get(); lastTestSucceededOk {
		model.LastTestSucceeded = types.BoolValue(lastTestSucceeded)
	}

	return model, diags
}
