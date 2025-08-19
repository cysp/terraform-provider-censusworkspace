package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

func SourceResourceIdentitySchema(_ context.Context) identityschema.Schema {
	return identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"id": identityschema.Int64Attribute{
				RequiredForImport: true,
			},
		},
	}
}

func SourceResourceSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"type": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				MarkdownDescription: "The type of the data source. A valid type is the `service_name` of a source type returned from the `/source_types` endpoint, where the source type is marked as `creatable_via_api`.",
			},
			"name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The name assigned to this source, typically a combination of type and location.",
			},
			"label": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "An optional label that can be assigned to the source for better categorization or identification.",
			},
			"sync_engine": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
				MarkdownDescription: "The sync engine type for this source. Can only be set during creation and cannot be modified after.",
			},
			"credentials": schema.StringAttribute{
				CustomType:          jsontypes.NormalizedType{},
				Optional:            true,
				MarkdownDescription: "Credentials that should be associated with this source (e.g. hostname, port)",
			},
			"created_at": schema.StringAttribute{
				CustomType:          timetypes.RFC3339Type{},
				Computed:            true,
				MarkdownDescription: "When the connection was created",
			},
			"last_tested_at": schema.StringAttribute{
				CustomType:          timetypes.RFC3339Type{},
				Computed:            true,
				MarkdownDescription: "Timestamp of when the last connection test was conducted on this source.",
			},
			"last_test_succeeded": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Indicates if the last connection test to this source was successful.",
			},
			"connection_details": schema.StringAttribute{
				CustomType:          jsontypes.NormalizedType{},
				Computed:            true,
				MarkdownDescription: "Detailed configuration and information for connecting to this source.",
			},
		},
	}
}
