package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

func SQLDatasetResourceIdentitySchema(_ context.Context) identityschema.Schema {
	return identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"id": identityschema.StringAttribute{
				RequiredForImport: true,
			},
		},
	}
}

func SQLDatasetResourceSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				MarkdownDescription: "Census identifier for this SQL dataset.",
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the dataset.",
			},
			"source_id": schema.Int64Attribute{
				Required: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
					int64planmodifier.UseStateForUnknown(),
				},
				MarkdownDescription: "ID of the source connection to run the query against.",
			},
			"query": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "SQL query that defines the dataset.",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Human-readable description shown for this SQL dataset in Census.",
			},
		},
	}
}
