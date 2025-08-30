package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type destinationModelBase struct {
	ID                types.String      `tfsdk:"id"`
	Name              types.String      `tfsdk:"name"`
	CreatedAt         timetypes.RFC3339 `tfsdk:"created_at"`
	LastTestedAt      timetypes.RFC3339 `tfsdk:"last_tested_at"`
	LastTestSucceeded types.Bool        `tfsdk:"last_test_succeeded"`
}

func destinationBaseResourceSchemaAttributes(_ context.Context) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"name": schema.StringAttribute{
			Required:            true,
			MarkdownDescription: "The name of this destination.",
		},
		"created_at": schema.StringAttribute{
			CustomType:          timetypes.RFC3339Type{},
			Computed:            true,
			MarkdownDescription: "When the connection was created",
		},
		"last_tested_at": schema.StringAttribute{
			CustomType:          timetypes.RFC3339Type{},
			Computed:            true,
			MarkdownDescription: "Timestamp of when the last connection test was conducted on this destination.",
		},
		"last_test_succeeded": schema.BoolAttribute{
			Computed:            true,
			MarkdownDescription: "Indicates if the last connection test to this destination was successful.",
		},
	}
}
