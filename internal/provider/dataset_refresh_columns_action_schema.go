package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/action/schema"
)

func DatasetRefreshColumnsActionSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"dataset_id": schema.Int64Attribute{
				Required: true,
			},
		},
	}
}
