package provider

import (
	"context"
	"maps"

	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func BigQuerySourceResourceIdentitySchema(_ context.Context) identityschema.Schema {
	return identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"id": identityschema.Int64Attribute{
				RequiredForImport: true,
			},
		},
	}
}

func BigQuerySourceResourceSchema(ctx context.Context) schema.Schema {
	attributes := sourceBaseResourceSchemaAttributes(ctx)

	maps.Copy(attributes, map[string]schema.Attribute{
		"credentials":        BigQuerySourceCredentialsResourceSchema(ctx),
		"connection_details": BigQuerySourceConnectionDetailsResourceSchema(ctx),
	})

	return schema.Schema{
		Attributes: attributes,
	}
}

//nolint:ireturn
func BigQuerySourceCredentialsResourceSchema(ctx context.Context) schema.Attribute {
	return schema.SingleNestedAttribute{
		Attributes: BigQuerySourceCredentialsResourceSchemaAttributes(ctx),
		CustomType: NewTypedObjectNull[BigQuerySourceCredentials]().CustomType(ctx),
		Required:   true,
	}
}

func BigQuerySourceCredentialsResourceSchemaAttributes(_ context.Context) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"project_id": schema.StringAttribute{
			Required: true,
		},
		"location": schema.StringAttribute{
			Required: true,
		},
		"service_account_key": schema.StringAttribute{
			Optional:  true,
			Sensitive: true,
		},
	}
}

//nolint:ireturn
func BigQuerySourceConnectionDetailsResourceSchema(ctx context.Context) schema.Attribute {
	return schema.SingleNestedAttribute{
		Attributes: BigQuerySourceConnectionDetailsResourceSchemaAttributes(ctx),
		CustomType: NewTypedObjectNull[BigQuerySourceConnectionDetails]().CustomType(ctx),
		Computed:   true,
	}
}

func BigQuerySourceConnectionDetailsResourceSchemaAttributes(_ context.Context) map[string]schema.Attribute {
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
