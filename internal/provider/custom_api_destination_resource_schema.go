package provider

import (
	"context"
	"maps"

	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
)

func CustomAPIDestinationResourceIdentitySchema(_ context.Context) identityschema.Schema {
	return identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"id": identityschema.Int64Attribute{
				RequiredForImport: true,
			},
		},
	}
}

func CustomAPIDestinationResourceSchema(ctx context.Context) schema.Schema {
	attributes := destinationBaseResourceSchemaAttributes(ctx)

	maps.Copy(attributes, map[string]schema.Attribute{
		"credentials":        CustomAPIDestinationCredentialsResourceSchema(ctx),
		"connection_details": CustomAPIDestinationConnectionDetailsResourceSchema(ctx),
	})

	return schema.Schema{
		Attributes: attributes,
	}
}

//nolint:ireturn
func CustomAPIDestinationCredentialsResourceSchema(ctx context.Context) schema.Attribute {
	return schema.SingleNestedAttribute{
		Attributes: CustomAPIDestinationCredentialsResourceSchemaAttributes(ctx),
		CustomType: NewTypedObjectNull[CustomAPIDestinationCredentials]().CustomType(ctx),
		Required:   true,
	}
}

func CustomAPIDestinationCredentialsResourceSchemaAttributes(ctx context.Context) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"api_version": schema.Int64Attribute{
			Required: true,
		},
		"webhook_url": schema.StringAttribute{
			Required: true,
		},
		"custom_headers": schema.MapNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: CustomAPIDestinationCredentialsCustomHeaderResourceSchemaAttributes(ctx),
				CustomType: NewTypedObjectNull[CustomAPIDestinationCustomHeader]().CustomType(ctx),
			},
			CustomType: NewTypedMapNull[TypedObject[CustomAPIDestinationCustomHeader]]().CustomType(ctx),
		},
	}
}

func CustomAPIDestinationCredentialsCustomHeaderResourceSchemaAttributes(_ context.Context) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"value": schema.StringAttribute{
			Required: true,
		},
		"is_secret": schema.BoolAttribute{
			Optional: true,
			Computed: true,
			Default:  booldefault.StaticBool(false),
		},
	}
}

//nolint:ireturn
func CustomAPIDestinationConnectionDetailsResourceSchema(ctx context.Context) schema.Attribute {
	return schema.SingleNestedAttribute{
		Attributes: CustomAPIDestinationConnectionDetailsResourceSchemaAttributes(ctx),
		CustomType: NewTypedObjectNull[CustomAPIDestinationConnectionDetails]().CustomType(ctx),
		Computed:   true,
	}
}

func CustomAPIDestinationConnectionDetailsResourceSchemaAttributes(ctx context.Context) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"api_version": schema.Int64Attribute{
			Computed: true,
		},
		"webhook_url": schema.StringAttribute{
			Computed: true,
		},
		"custom_headers": schema.MapNestedAttribute{
			Computed: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: CustomAPIDestinationConnectionDetailsCustomHeaderResourceSchemaAttributes(ctx),
				CustomType: NewTypedObjectNull[CustomAPIDestinationCustomHeader]().CustomType(ctx),
			},
			CustomType: NewTypedMapNull[TypedObject[CustomAPIDestinationCustomHeader]]().CustomType(ctx),
		},
	}
}

func CustomAPIDestinationConnectionDetailsCustomHeaderResourceSchemaAttributes(_ context.Context) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"value": schema.StringAttribute{
			Computed: true,
		},
		"is_secret": schema.BoolAttribute{
			Computed: true,
		},
	}
}
