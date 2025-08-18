package provider

import (
	"context"
	"fmt"

	cm "github.com/cysp/terraform-provider-censusworkspace/internal/census-management-go"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewBigQuerySourceModelFromResponse(ctx context.Context, response cm.SourceData) (BigQuerySourceModel, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	path := path.Empty()

	model := BigQuerySourceModel{
		sourceModelCommon: sourceModelCommon{
			ID:        types.StringValue(fmt.Sprintf("%d", response.ID)),
			Name:      types.StringValue(response.Name),
			Label:     types.StringPointerValue(response.Label.ValueStringPointer()),
			CreatedAt: timetypes.NewRFC3339TimeValue(response.CreatedAt),
		},
	}

	if response.ConnectionDetails != nil {
		path := path.AtName("connection_details")

		connectionDetails, connectionDetailsDiags := NewBigQuerySourceConnectionDetailsFromResponse(ctx, path, response.ConnectionDetails)
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

// func (s *BigQuerySourceConnectionDetails) Decode(e *jx.Decoder) error {
// 	return e.Obj(func(e *jx.Decoder, key string) error {
// 		switch key {
// 		case "project_id":
// 			value, err := e.Str()
// 			if err != nil {
// 				return err
// 			}

// 			s.ProjectID = types.StringValue(value)

// 		case "location":
// 			value, err := e.Str()
// 			if err != nil {
// 				return err
// 			}

// 			s.Location = types.StringValue(value)

// 		case "service_account":
// 			value, err := e.Str()
// 			if err != nil {
// 				return err
// 			}

// 			s.ServiceAccount = types.StringValue(value)
// 		}

// 		return nil
// 	})
// }
