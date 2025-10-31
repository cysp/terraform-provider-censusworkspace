package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type SyncModel struct {
	ID        types.String `tfsdk:"id"`
	Label     types.String `tfsdk:"name"`
	Operation types.String `tfsdk:"operation"`
	// SourceID    types.Int64       `tfsdk:"source_id"`
	// Description types.String      `tfsdk:"description"`
	// CreatedAt   timetypes.RFC3339 `tfsdk:"created_at"`
	// UpdatedAt   timetypes.RFC3339 `tfsdk:"updated_at"`
}
