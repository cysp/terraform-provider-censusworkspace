package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type sourceModelBase struct {
	ID                                types.String `tfsdk:"id"`
	Name                              types.String `tfsdk:"name"`
	Label                             types.String `tfsdk:"label"`
	SyncEngine                        types.String `tfsdk:"sync_engine"`
	WarehouseWritebackRetentionInDays types.Int64  `tfsdk:"warehouse_writeback_retention_in_days"`
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
			Required:            true,
			MarkdownDescription: "The name assigned to this source.",
		},
		"label": schema.StringAttribute{
			Computed:            true,
			DeprecationMessage:  "Use name instead.",
			MarkdownDescription: "Deprecated. Use `name` for configuration. This read-only field reflects the API label when returned.",
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
		"warehouse_writeback_retention_in_days": schema.Int64Attribute{
			Optional: true,
			Computed: true,
			Validators: []validator.Int64{
				int64validator.AtLeast(1),
			},
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
			MarkdownDescription: "Number of days to retain warehouse writeback data. When set, automatically enables sync logs for this source.",
		},
	}
}
