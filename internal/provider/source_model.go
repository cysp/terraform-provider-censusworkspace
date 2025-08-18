package provider

import (
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type SourceModel struct {
	sourceModelCommon

	Credentials       jsontypes.Normalized `tfsdk:"credentials"`
	ConnectionDetails jsontypes.Normalized `tfsdk:"connection_details"`
}

type sourceModelCommon struct {
	ID                types.String      `tfsdk:"id"`
	Name              types.String      `tfsdk:"name"`
	Type              types.String      `tfsdk:"type"`
	SyncEngine        types.String      `tfsdk:"sync_engine"`
	Label             types.String      `tfsdk:"label"`
	CreatedAt         timetypes.RFC3339 `tfsdk:"created_at"`
	LastTestedAt      timetypes.RFC3339 `tfsdk:"last_tested_at"`
	LastTestSucceeded types.Bool        `tfsdk:"last_test_succeeded"`
}
