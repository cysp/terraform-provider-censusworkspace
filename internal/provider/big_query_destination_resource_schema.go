package provider

import (
	"context"
	"maps"

	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func BigQueryDestinationResourceIdentitySchema(_ context.Context) identityschema.Schema {
	return identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"id": identityschema.StringAttribute{
				RequiredForImport: true,
			},
		},
	}
}

func BigQueryDestinationResourceSchema(ctx context.Context) schema.Schema {
	attributes := destinationBaseResourceSchemaAttributes(ctx)

	maps.Copy(attributes, map[string]schema.Attribute{
		"credentials":        BigQueryDestinationCredentialsResourceSchema(ctx),
		"connection_details": BigQueryDestinationConnectionDetailsResourceSchema(ctx),
	})

	return schema.Schema{
		Attributes: attributes,
	}
}

//nolint:ireturn
func BigQueryDestinationCredentialsResourceSchema(ctx context.Context) schema.Attribute {
	return schema.SingleNestedAttribute{
		Attributes:          BigQueryDestinationCredentialsResourceSchemaAttributes(ctx),
		CustomType:          NewTypedObjectNull[BigQueryDestinationCredentials]().CustomType(ctx),
		Required:            true,
		MarkdownDescription: "BigQuery destination connection values to send to Census.",
	}
}

func BigQueryDestinationCredentialsResourceSchemaAttributes(_ context.Context) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"project_id": schema.StringAttribute{
			Required:            true,
			MarkdownDescription: "Google Cloud project ID containing the BigQuery destination dataset.",
		},
		"location": schema.StringAttribute{
			Required:            true,
			MarkdownDescription: "BigQuery location of the destination dataset.",
		},
		"service_account_key": schema.StringAttribute{
			Optional:            true,
			Sensitive:           true,
			MarkdownDescription: "Service account key JSON used by Census to write to BigQuery.",
		},
	}
}

//nolint:ireturn
func BigQueryDestinationConnectionDetailsResourceSchema(ctx context.Context) schema.Attribute {
	return schema.SingleNestedAttribute{
		Attributes:          BigQueryDestinationConnectionDetailsResourceSchemaAttributes(ctx),
		CustomType:          NewTypedObjectNull[BigQueryDestinationConnectionDetails]().CustomType(ctx),
		Computed:            true,
		MarkdownDescription: "BigQuery destination connection details returned by Census.",
	}
}

func BigQueryDestinationConnectionDetailsResourceSchemaAttributes(_ context.Context) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"project_id": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "Google Cloud project ID Census has stored for this BigQuery destination.",
		},
		"location": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "BigQuery location Census has stored for this BigQuery destination.",
		},
		"service_account_email": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "Service account email Census uses for this BigQuery destination.",
		},
		"service_account_key": schema.StringAttribute{
			Computed:            true,
			Sensitive:           true,
			MarkdownDescription: "Service account key value returned by Census for this BigQuery destination, when available.",
		},
	}
}
