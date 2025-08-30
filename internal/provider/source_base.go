package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type sourceModelBase struct {
	ID                types.String      `tfsdk:"id"`
	Name              types.String      `tfsdk:"name"`
	Label             types.String      `tfsdk:"label"`
	SyncEngine        types.String      `tfsdk:"sync_engine"`
	CreatedAt         timetypes.RFC3339 `tfsdk:"created_at"`
	LastTestedAt      timetypes.RFC3339 `tfsdk:"last_tested_at"`
	LastTestSucceeded types.Bool        `tfsdk:"last_test_succeeded"`
}

func sourceBaseResourceSchemaAttributes(_ context.Context) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"name": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "The name assigned to this source, typically a combination of type and location.",
		},
		"label": schema.StringAttribute{
			Required:            true,
			MarkdownDescription: "An optional label that can be assigned to the source for better categorization or identification.",
		},
		"sync_engine": schema.StringAttribute{
			Optional: true,
			Computed: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
				stringplanmodifier.UseStateForUnknown(),
			},
			MarkdownDescription: "The sync engine type for this source. Can only be set during creation and cannot be modified after.",
		},
		"created_at": schema.StringAttribute{
			CustomType:          timetypes.RFC3339Type{},
			Computed:            true,
			MarkdownDescription: "When the connection was created",
		},
		"last_tested_at": schema.StringAttribute{
			CustomType:          timetypes.RFC3339Type{},
			Computed:            true,
			MarkdownDescription: "Timestamp of when the last connection test was conducted on this source.",
		},
		"last_test_succeeded": schema.BoolAttribute{
			Computed:            true,
			MarkdownDescription: "Indicates if the last connection test to this source was successful.",
		},
	}
}
