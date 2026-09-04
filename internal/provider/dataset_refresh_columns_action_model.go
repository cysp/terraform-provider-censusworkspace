package provider

import "github.com/hashicorp/terraform-plugin-framework/types"

type DatasetRefreshColumnsActionModel struct {
	DatasetID types.Int64 `tfsdk:"dataset_id"`
}
