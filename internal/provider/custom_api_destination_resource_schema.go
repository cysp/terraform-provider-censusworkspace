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
			"id": identityschema.StringAttribute{
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
		Attributes:          CustomAPIDestinationCredentialsResourceSchemaAttributes(ctx),
		CustomType:          NewTypedObjectNull[CustomAPIDestinationCredentials]().CustomType(ctx),
		Required:            true,
		MarkdownDescription: "Custom API destination connection values to send to Census.",
	}
}

func CustomAPIDestinationCredentialsResourceSchemaAttributes(ctx context.Context) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"api_version": schema.Int64Attribute{
			Required:            true,
			MarkdownDescription: "Custom Destination API version implemented by your endpoint.",
		},
		"webhook_url": schema.StringAttribute{
			Required:            true,
			MarkdownDescription: "Public webhook URL Census calls for Custom API destination syncs.",
		},
		"custom_headers": schema.MapNestedAttribute{
			Optional:            true,
			MarkdownDescription: "Additional HTTP headers Census sends to the Custom API webhook.",
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
			Optional:            true,
			MarkdownDescription: "HTTP header value sent to the Custom API webhook. Configure exactly one of `value` or `value_wo` for each header.",
		},
		"value_wo": schema.StringAttribute{
			Optional:            true,
			Sensitive:           true,
			WriteOnly:           true,
			MarkdownDescription: "HTTP header value sent to the Custom API webhook. Configure exactly one of `value` or `value_wo` for each header.",
		},
		"is_secret": schema.BoolAttribute{
			Optional:            true,
			Computed:            true,
			Default:             booldefault.StaticBool(false),
			MarkdownDescription: "Whether Census should store this header value as a secret.",
		},
	}
}

//nolint:ireturn
func CustomAPIDestinationConnectionDetailsResourceSchema(ctx context.Context) schema.Attribute {
	return schema.SingleNestedAttribute{
		Attributes:          CustomAPIDestinationConnectionDetailsResourceSchemaAttributes(ctx),
		CustomType:          NewTypedObjectNull[CustomAPIDestinationConnectionDetails]().CustomType(ctx),
		Computed:            true,
		MarkdownDescription: "Custom API destination connection details returned by Census.",
	}
}

func CustomAPIDestinationConnectionDetailsResourceSchemaAttributes(ctx context.Context) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"api_version": schema.Int64Attribute{
			Computed:            true,
			MarkdownDescription: "Custom Destination API version Census has stored for this destination.",
		},
		"webhook_url": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "Public webhook URL Census has stored for this Custom API destination.",
		},
		"custom_headers": schema.MapNestedAttribute{
			Computed:            true,
			MarkdownDescription: "Additional HTTP headers Census has stored for this Custom API webhook.",
			NestedObject: schema.NestedAttributeObject{
				Attributes: CustomAPIDestinationConnectionDetailsCustomHeaderResourceSchemaAttributes(ctx),
				CustomType: NewTypedObjectNull[CustomAPIDestinationConnectionDetailsCustomHeader]().CustomType(ctx),
			},
			CustomType: NewTypedMapNull[TypedObject[CustomAPIDestinationConnectionDetailsCustomHeader]]().CustomType(ctx),
		},
	}
}

func CustomAPIDestinationConnectionDetailsCustomHeaderResourceSchemaAttributes(_ context.Context) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"value": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "HTTP header value returned by Census, or null when the header is secret.",
		},
		"is_secret": schema.BoolAttribute{
			Computed:            true,
			MarkdownDescription: "Whether Census stores this header value as a secret.",
		},
	}
}
