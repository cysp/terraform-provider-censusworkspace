package provider

import (
	"context"
	"maps"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func BigQuerySourceResourceIdentitySchema(_ context.Context) identityschema.Schema {
	return identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"id": identityschema.StringAttribute{
				RequiredForImport: true,
			},
		},
	}
}

func BigQuerySourceResourceSchema(ctx context.Context) schema.Schema {
	attributes := sourceBaseResourceSchemaAttributes(ctx)

	maps.Copy(attributes, map[string]schema.Attribute{
		"credentials":        BigQuerySourceCredentialsResourceSchema(ctx),
		"connection_details": BigQuerySourceConnectionDetailsResourceSchema(ctx),
	})

	return schema.Schema{
		Attributes: attributes,
	}
}

//nolint:ireturn
func BigQuerySourceCredentialsResourceSchema(ctx context.Context) schema.Attribute {
	return schema.SingleNestedAttribute{
		Attributes:          BigQuerySourceCredentialsResourceSchemaAttributes(ctx),
		CustomType:          NewTypedObjectNull[BigQuerySourceCredentials]().CustomType(ctx),
		Required:            true,
		MarkdownDescription: "BigQuery source connection values to send to Census.",
	}
}

func BigQuerySourceCredentialsResourceSchemaAttributes(ctx context.Context) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"project_id": schema.StringAttribute{
			Required:            true,
			MarkdownDescription: "Google Cloud project ID containing the BigQuery source data.",
		},
		"location": schema.StringAttribute{
			Required:            true,
			MarkdownDescription: "BigQuery location of the source data.",
		},
		"service_account_key": schema.SingleNestedAttribute{
			Optional:            true,
			Attributes:          BigQuerySourceCredentialsServiceAccountKeyResourceSchemaAttributes(ctx),
			CustomType:          NewTypedObjectNull[BigQuerySourceCredentialsServiceAccountKey]().CustomType(ctx),
			MarkdownDescription: "Service account key JSON fields used by Census to read from BigQuery.",
		},
	}
}

func BigQuerySourceCredentialsServiceAccountKeyResourceSchemaAttributes(_ context.Context) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"type": schema.StringAttribute{
			Required: true,
			Validators: []validator.String{
				stringvalidator.OneOf("service_account"),
			},
			MarkdownDescription: "Service account key type. Must be `service_account`.",
		},
		"project_id": schema.StringAttribute{
			Required:            true,
			MarkdownDescription: "Google Cloud project ID from the service account key JSON.",
		},
		"private_key_id": schema.StringAttribute{
			Required:            true,
			MarkdownDescription: "Private key ID from the service account key JSON.",
		},
		"private_key": schema.StringAttribute{
			Sensitive:           true,
			Required:            true,
			MarkdownDescription: "Private key from the service account key JSON.",
		},
		"client_email": schema.StringAttribute{
			Required:            true,
			MarkdownDescription: "Client email from the service account key JSON.",
		},
		"client_id": schema.StringAttribute{
			Required:            true,
			MarkdownDescription: "Client ID from the service account key JSON.",
		},
	}
}

//nolint:ireturn
func BigQuerySourceConnectionDetailsResourceSchema(ctx context.Context) schema.Attribute {
	return schema.SingleNestedAttribute{
		Attributes:          BigQuerySourceConnectionDetailsResourceSchemaAttributes(ctx),
		CustomType:          NewTypedObjectNull[BigQuerySourceConnectionDetails]().CustomType(ctx),
		Computed:            true,
		MarkdownDescription: "BigQuery source connection details returned by Census.",
	}
}

func BigQuerySourceConnectionDetailsResourceSchemaAttributes(_ context.Context) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"project_id": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "Google Cloud project ID Census has stored for this BigQuery source.",
		},
		"location": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "BigQuery location Census has stored for this BigQuery source.",
		},
		"service_account": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "Census-managed service account email address for this BigQuery source.",
		},
	}
}
