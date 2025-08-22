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

func DestinationResourceIdentitySchema(_ context.Context) identityschema.Schema {
	return identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"id": identityschema.Int64Attribute{
				RequiredForImport: true,
			},
		},
	}
}

func DestinationResourceSchema(_ context.Context) schema.Schema {
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
				MarkdownDescription: "The type of connection to be used for this destination. A valid type is the `service_name` of a connector returned from the `/connectors` endpoint, where the connector is marked as `creatable_via_api`.",
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of this destination.",
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
			"created_at": schema.StringAttribute{
				CustomType:          timetypes.RFC3339Type{},
				Computed:            true,
				MarkdownDescription: "When the connection was created",
			},
			"last_tested_at": schema.StringAttribute{
				CustomType:          timetypes.RFC3339Type{},
				Computed:            true,
				MarkdownDescription: "Timestamp of when the last connection test was conducted on this destination.",
			},
			"last_test_succeeded": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Indicates if the last connection test to this destination was successful.",
			},
		},
	}
}
