package provider

import (
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type SourceModel struct {
	ID                types.String         `tfsdk:"id"`
	Name              types.String         `tfsdk:"name"`
	Type              types.String         `tfsdk:"type"`
	SyncEngine        types.String         `tfsdk:"sync_engine"`
	Label             types.String         `tfsdk:"label"`
	Credentials       jsontypes.Normalized `tfsdk:"credentials"`
	ConnectionDetails jsontypes.Normalized `tfsdk:"connection_details"`
	CreatedAt         timetypes.RFC3339    `tfsdk:"created_at"`
	LastTestedAt      timetypes.RFC3339    `tfsdk:"last_tested_at"`
	LastTestSucceeded types.Bool           `tfsdk:"last_test_succeeded"`
}
