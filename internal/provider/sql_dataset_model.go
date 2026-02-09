package provider

import "github.com/hashicorp/terraform-plugin-framework/types"

type SQLDatasetModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	SourceID    types.Int64  `tfsdk:"source_id"`
	Query       types.String `tfsdk:"query"`
	Description types.String `tfsdk:"description"`
}
