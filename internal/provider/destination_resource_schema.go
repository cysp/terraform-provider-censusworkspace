package provider

import (
	"context"
	"maps"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

func DestinationResourceIdentitySchema(_ context.Context) identityschema.Schema {
	return identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"id": identityschema.Int64Attribute{
				RequiredForImport: true,
			},
		},
	}
}

func DestinationResourceSchema(ctx context.Context) schema.Schema {
	attributes := destinationBaseResourceSchemaAttributes(ctx)

	maps.Copy(attributes, map[string]schema.Attribute{
		"type": schema.StringAttribute{
			Required: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
			MarkdownDescription: "The type of connection to be used for this destination. A valid type is the `service_name` of a connector returned from the `/connectors` endpoint, where the connector is marked as `creatable_via_api`.",
		},
		"credentials": schema.StringAttribute{
			CustomType:          jsontypes.NormalizedType{},
			Optional:            true,
			MarkdownDescription: "The credentials needed to create each type of connection. These can be found in the `GET /connectors` API for most destinations.",
		},
		"connection_details": schema.StringAttribute{
			CustomType:          jsontypes.NormalizedType{},
			Computed:            true,
			MarkdownDescription: "Connection details associated with this destination.",
		},
	})

	return schema.Schema{
		Attributes: attributes,
	}
}
