package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

func SourceResourceSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					UseStateForUnknown(),
				},
			},
			"type": schema.StringAttribute{
				Required: true,
			},
			"name": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					UseStateForUnknown(),
				},
			},
			"label": schema.StringAttribute{
				Optional: true,
			},
			"credentials": schema.StringAttribute{
				CustomType: jsontypes.NormalizedType{},
				Optional:   true,
			},
			"created_at": schema.StringAttribute{
				CustomType: timetypes.RFC3339Type{},
				Computed:   true,
				PlanModifiers: []planmodifier.String{
					UseStateForUnknown(),
				},
			},
			"last_tested_at": schema.StringAttribute{
				CustomType: timetypes.RFC3339Type{},
				Computed:   true,
			},
			"last_test_succeeded": schema.BoolAttribute{
				Computed: true,
			},
			"connection_details": schema.StringAttribute{
				CustomType: jsontypes.NormalizedType{},
				Computed:   true,
			},
		},
	}
}
