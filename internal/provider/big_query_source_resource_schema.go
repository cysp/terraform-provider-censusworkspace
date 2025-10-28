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
		Attributes: BigQuerySourceCredentialsResourceSchemaAttributes(ctx),
		CustomType: NewTypedObjectNull[BigQuerySourceCredentials]().CustomType(ctx),
		Required:   true,
	}
}

func BigQuerySourceCredentialsResourceSchemaAttributes(ctx context.Context) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"project_id": schema.StringAttribute{
			Required: true,
		},
		"location": schema.StringAttribute{
			Required: true,
		},
		"service_account_key": schema.SingleNestedAttribute{
			Optional:   true,
			Attributes: BigQuerySourceCredentialsServiceAccountKeyResourceSchemaAttributes(ctx),
			CustomType: NewTypedObjectNull[BigQuerySourceCredentialsServiceAccountKey]().CustomType(ctx),
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
		},
		"project_id": schema.StringAttribute{
			Required: true,
		},
		"private_key_id": schema.StringAttribute{
			Required: true,
		},
		"private_key": schema.StringAttribute{
			Sensitive: true,
			Required:  true,
		},
		"client_email": schema.StringAttribute{
			Required: true,
		},
		"client_id": schema.StringAttribute{
			Required: true,
		},
	}
}

//nolint:ireturn
func BigQuerySourceConnectionDetailsResourceSchema(ctx context.Context) schema.Attribute {
	return schema.SingleNestedAttribute{
		Attributes: BigQuerySourceConnectionDetailsResourceSchemaAttributes(ctx),
		CustomType: NewTypedObjectNull[BigQuerySourceConnectionDetails]().CustomType(ctx),
		Computed:   true,
	}
}

func BigQuerySourceConnectionDetailsResourceSchemaAttributes(_ context.Context) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"project_id": schema.StringAttribute{
			Computed: true,
		},
		"location": schema.StringAttribute{
			Computed: true,
		},
		"service_account": schema.StringAttribute{
			Computed: true,
		},
	}
}
