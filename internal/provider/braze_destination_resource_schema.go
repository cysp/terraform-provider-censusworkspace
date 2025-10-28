package provider

import (
	"context"
	"maps"

	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func BrazeDestinationResourceIdentitySchema(_ context.Context) identityschema.Schema {
	return identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"id": identityschema.StringAttribute{
				RequiredForImport: true,
			},
		},
	}
}

func BrazeDestinationResourceSchema(ctx context.Context) schema.Schema {
	attributes := destinationBaseResourceSchemaAttributes(ctx)

	maps.Copy(attributes, map[string]schema.Attribute{
		"credentials":        BrazeDestinationCredentialsResourceSchema(ctx),
		"connection_details": BrazeDestinationConnectionDetailsResourceSchema(ctx),
	})

	return schema.Schema{
		Attributes: attributes,
	}
}

//nolint:ireturn
func BrazeDestinationCredentialsResourceSchema(ctx context.Context) schema.Attribute {
	return schema.SingleNestedAttribute{
		Attributes: BrazeDestinationCredentialsResourceSchemaAttributes(ctx),
		CustomType: NewTypedObjectNull[BrazeDestinationCredentials]().CustomType(ctx),
		Required:   true,
	}
}

func BrazeDestinationCredentialsResourceSchemaAttributes(_ context.Context) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"instance_url": schema.StringAttribute{
			Required:    true,
			Description: "Endpoint URL",
		},
		"api_key": schema.StringAttribute{
			Required:    true,
			Sensitive:   true,
			Description: "API Key",
		},
		"client_key": schema.StringAttribute{
			Optional:    true,
			Sensitive:   true,
			Description: "Data Import Key (for Cohorts only)",
		},
	}
}

//nolint:ireturn
func BrazeDestinationConnectionDetailsResourceSchema(ctx context.Context) schema.Attribute {
	return schema.SingleNestedAttribute{
		Attributes: BrazeDestinationConnectionDetailsResourceSchemaAttributes(ctx),
		CustomType: NewTypedObjectNull[BrazeDestinationConnectionDetails]().CustomType(ctx),
		Computed:   true,
	}
}

func BrazeDestinationConnectionDetailsResourceSchemaAttributes(_ context.Context) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"instance_url": schema.StringAttribute{
			Computed:    true,
			Description: "Endpoint URL",
		},
	}
}
