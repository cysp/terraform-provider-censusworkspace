package provider

import (
	"context"
	"strconv"

	cm "github.com/cysp/terraform-provider-censusworkspace/internal/census-management-go"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewDestinationModelFromResponse(_ context.Context, response cm.DestinationData) (DestinationModel, diag.Diagnostics) {
	model := DestinationModel{
		destinationModelBase: destinationModelBase{
			ID:        types.StringValue(strconv.FormatInt(response.ID, 10)),
			Name:      types.StringValue(response.Name),
			CreatedAt: timetypes.NewRFC3339TimeValue(response.CreatedAt),
		},
		Type: types.StringValue(response.Type),
	}

	if response.ConnectionDetails != nil {
		model.ConnectionDetails = jsontypes.NewNormalizedValue(string(response.ConnectionDetails))
	}

	if lastTestedAt, lastTestedAtOk := response.LastTestedAt.Get(); lastTestedAtOk {
		model.LastTestedAt = timetypes.NewRFC3339TimeValue(lastTestedAt)
	}

	if lastTestSucceeded, lastTestSucceededOk := response.LastTestSucceeded.Get(); lastTestSucceededOk {
		model.LastTestSucceeded = types.BoolValue(lastTestSucceeded)
	}

	return model, nil
}
