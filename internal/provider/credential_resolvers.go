package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type stringCredentialInput struct {
	Legacy        types.String
	WriteOnly     types.String
	LegacyPath    path.Path
	WriteOnlyPath path.Path
	Required      bool
}

type resolvedStringCredential struct {
	Value         types.String
	UsedWriteOnly bool
	WriteOnlyPath path.Path
}

func (c resolvedStringCredential) AddValueTo(values writeOnlyCredentialValues) diag.Diagnostics {
	if c.UsedWriteOnly {
		return values.Add(c.WriteOnlyPath, c.Value)
	}

	return nil
}

func resolveStringCredentialInput(input stringCredentialInput) (resolvedStringCredential, diag.Diagnostics) {
	value, usedWriteOnly, diags := resolveStringCredential(input.Legacy, input.WriteOnly, input.LegacyPath, input.WriteOnlyPath, input.Required)

	return resolvedStringCredential{
		Value:         value,
		UsedWriteOnly: usedWriteOnly,
		WriteOnlyPath: input.WriteOnlyPath,
	}, diags
}

func configuredObjectOrPlan[T any](configured, planned TypedObject[T]) TypedObject[T] {
	if configured.IsNull() || configured.IsUnknown() {
		return planned
	}

	return configured
}
