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

func SourceResourceIdentitySchema(_ context.Context) identityschema.Schema {
	return identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"id": identityschema.StringAttribute{
				RequiredForImport: true,
			},
		},
	}
}

func SourceResourceSchema(ctx context.Context) schema.Schema {
	attributes := sourceBaseResourceSchemaAttributes(ctx)

	maps.Copy(attributes, map[string]schema.Attribute{
		"type": schema.StringAttribute{
			Required: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
			MarkdownDescription: "The type of the data source. A valid type is the `service_name` of a source type returned from the `/source_types` endpoint, where the source type is marked as `creatable_via_api`.",
		},
		"credentials": schema.StringAttribute{
			CustomType:          jsontypes.NormalizedType{},
			Optional:            true,
			MarkdownDescription: "Credentials that should be associated with this source (e.g. hostname, port)",
		},
		"connection_details": schema.StringAttribute{
			CustomType:          jsontypes.NormalizedType{},
			Computed:            true,
			MarkdownDescription: "Detailed configuration and information for connecting to this source.",
		},
	})

	return schema.Schema{
		Attributes: attributes,
	}
}
