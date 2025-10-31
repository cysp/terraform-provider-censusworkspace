package provider

import (
	"context"
	"maps"

	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

func SyncResourceIdentitySchema(_ context.Context) identityschema.Schema {
	return identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"id": identityschema.StringAttribute{
				RequiredForImport: true,
			},
		},
	}
}

func SyncResourceSchema(ctx context.Context) schema.Schema {
	attributes := syncBaseResourceSchemaAttributes(ctx)

	maps.Copy(attributes, map[string]schema.Attribute{
		"operation": schema.StringAttribute{
			Required:            true,
			MarkdownDescription: "How records are synced to the destination.",
		},
		// "credentials": schema.StringAttribute{
		// 	CustomType:          jsontypes.NormalizedType{},
		// 	Optional:            true,
		// 	MarkdownDescription: "The credentials needed to create each type of connection. These can be found in the `GET /connectors` API for most syncs.",
		// },
		// "connection_details": schema.StringAttribute{
		// 	CustomType:          jsontypes.NormalizedType{},
		// 	Computed:            true,
		// 	MarkdownDescription: "Connection details associated with this sync.",
		// },
	})

	return schema.Schema{
		Attributes: attributes,
	}
}

func syncBaseResourceSchemaAttributes(_ context.Context) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"label": schema.StringAttribute{
			Required:            true,
			MarkdownDescription: "A label to give to this sync.",
		},
	}
}
