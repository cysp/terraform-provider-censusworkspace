package provider

import (
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type SourceModel struct {
	ID                types.String         `tfsdk:"id"`
	Type              types.String         `tfsdk:"type"`
	Name              types.String         `tfsdk:"name"`
	Label             types.String         `tfsdk:"label"`
	Credentials       jsontypes.Normalized `tfsdk:"credentials"`
	CreatedAt         timetypes.RFC3339    `tfsdk:"created_at"`
	LastTestedAt      timetypes.RFC3339    `tfsdk:"last_tested_at"`
	LastTestSucceeded types.Bool           `tfsdk:"last_test_succeeded"`
	ConnectionDetails jsontypes.Normalized `tfsdk:"connection_details"`
}

// The unique identifier of the source.
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
