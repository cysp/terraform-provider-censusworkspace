package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

// ID          types.Int64          `tfsdk:"id"`
// Type        types.String         `tfsdk:"type"`
// Label       types.String         `tfsdk:"label"`
// Credentials jsontypes.Normalized `tfsdk:"credentials"`

// // The unique identifier of the source.
// ID int64 `json:"id"`
// // The name assigned to this source, typically a combination of type and location.
// Name string `json:"name"`
// // An optional label that can be assigned to the source for better categorization or identification.
// Label OptNilString `json:"label"`
// // The type of the data source. A valid type is the service_name of a source type returned from the
// // /source_types endpoint, where the source type is marked as creatable_via_api.
// Type string `json:"type"`
// // The timestamp when the source was created.
// CreatedAt time.Time `json:"created_at"`
// // Indicates if the last connection test to this source was successful.
// LastTestSucceeded OptNilBool `json:"last_test_succeeded"`
// // Timestamp of when the last connection test was conducted on this source.
// LastTestedAt      OptNilDateTime `json:"last_tested_at"`
// ConnectionDetails jx.Raw         `json:"connection_details"`
func SourceResourceSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					UseStateForUnknown(),
				},
			},
			"type": schema.StringAttribute{
				Required: true,
			},
			"name": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					UseStateForUnknown(),
				},
			},
			"label": schema.StringAttribute{
				Optional: true,
			},
			"credentials": schema.StringAttribute{
				CustomType: jsontypes.NormalizedType{},
				Optional:   true,
			},
			"created_at": schema.StringAttribute{
				CustomType: timetypes.RFC3339Type{},
				Computed:   true,
				PlanModifiers: []planmodifier.String{
					UseStateForUnknown(),
				},
			},
			"last_tested_at": schema.StringAttribute{
				CustomType: timetypes.RFC3339Type{},
				Computed:   true,
			},
			"last_test_succeeded": schema.BoolAttribute{
				Computed: true,
			},
			"connection_details": schema.StringAttribute{
				CustomType: jsontypes.NormalizedType{},
				Computed:   true,
			},
		},
	}
}
