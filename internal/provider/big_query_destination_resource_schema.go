package provider

import (
	"context"
	"maps"

	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func BigQueryDestinationResourceIdentitySchema(_ context.Context) identityschema.Schema {
	return identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"id": identityschema.Int64Attribute{
				RequiredForImport: true,
			},
		},
	}
}

func BigQueryDestinationResourceSchema(ctx context.Context) schema.Schema {
	attributes := destinationBaseResourceSchemaAttributes(ctx)

	maps.Copy(attributes, map[string]schema.Attribute{
		"credentials":        BigQueryDestinationCredentialsResourceSchema(ctx),
		"connection_details": BigQueryDestinationConnectionDetailsResourceSchema(ctx),
	})

	return schema.Schema{
		Attributes: attributes,
	}
}

//nolint:ireturn
func BigQueryDestinationCredentialsResourceSchema(ctx context.Context) schema.Attribute {
	return schema.SingleNestedAttribute{
		Attributes: BigQueryDestinationCredentialsResourceSchemaAttributes(ctx),
		CustomType: NewTypedObjectNull[BigQueryDestinationCredentials]().CustomType(ctx),
		Required:   true,
	}
}

func BigQueryDestinationCredentialsResourceSchemaAttributes(_ context.Context) map[string]schema.Attribute {
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
func BigQueryDestinationConnectionDetailsResourceSchema(ctx context.Context) schema.Attribute {
	return schema.SingleNestedAttribute{
		Attributes: BigQueryDestinationConnectionDetailsResourceSchemaAttributes(ctx),
		CustomType: NewTypedObjectNull[BigQueryDestinationConnectionDetails]().CustomType(ctx),
		Computed:   true,
	}
}

func BigQueryDestinationConnectionDetailsResourceSchemaAttributes(_ context.Context) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"project_id": schema.StringAttribute{
			Computed: true,
		},
		"location": schema.StringAttribute{
			Computed: true,
		},
		"service_account_email": schema.StringAttribute{
			Computed: true,
		},
		"service_account_key": schema.StringAttribute{
			Computed:  true,
			Sensitive: true,
		},
	}
}
