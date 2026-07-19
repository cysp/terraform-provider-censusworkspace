package provider

import (
	"context"

	"github.com/cysp/terraform-provider-censusworkspace/internal/census"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ validator.String = canonicalSourceNameValidator{}

type canonicalSourceNameValidator struct{}

func (canonicalSourceNameValidator) Description(context.Context) string {
	return "Source names must use underscores instead of spaces or hyphens."
}

func (v canonicalSourceNameValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v canonicalSourceNameValidator) ValidateString(
	ctx context.Context,
	request validator.StringRequest,
	response *validator.StringResponse,
) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	name := request.ConfigValue.ValueString()
	if census.CanonicalizeSourceName(name) == name {
		return
	}

	response.Diagnostics.AddAttributeError(
		request.Path,
		"Noncanonical Source Name",
		v.Description(ctx),
	)
}
