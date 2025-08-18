package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

//nolint:ireturn
func BigQuerySourceConnectionDetailsSchema(ctx context.Context) schema.Attribute {
	return schema.SingleNestedAttribute{
		Attributes: BigQuerySourceConnectionDetails{}.SchemaAttributes(ctx),
		CustomType: BigQuerySourceConnectionDetails{}.CustomType(ctx),
		Computed:   true,
	}
}

func (v BigQuerySourceConnectionDetails) SchemaAttributes(_ context.Context) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"project_id": schema.StringAttribute{
			Computed: true,
		},
		"location": schema.StringAttribute{
			Computed: true,
		},
		"service_account": schema.StringAttribute{
			Computed: true,
		},
	}
}
