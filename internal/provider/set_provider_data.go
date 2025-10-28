package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func SetProviderDataFromActionConfigureRequest[ProviderData any](req action.ConfigureRequest, out *ProviderData) diag.Diagnostics {
	diags := diag.Diagnostics{}

	if req.ProviderData == nil {
		return diags
	}

	if providerData, ok := req.ProviderData.(ProviderData); ok {
		*out = providerData

		return diags
	}

	diags.AddError("Invalid provider data", "")

	return diags
}

func SetProviderDataFromDataSourceConfigureRequest[ProviderData any](req datasource.ConfigureRequest, out *ProviderData) diag.Diagnostics {
	diags := diag.Diagnostics{}

	if req.ProviderData == nil {
		return diags
	}

	if providerData, ok := req.ProviderData.(ProviderData); ok {
		*out = providerData

		return diags
	}

	diags.AddError("Invalid provider data", "")

	return diags
}

func SetProviderDataFromResourceConfigureRequest[ProviderData any](req resource.ConfigureRequest, out *ProviderData) diag.Diagnostics {
	diags := diag.Diagnostics{}

	if req.ProviderData == nil {
		return diags
	}

	if providerData, ok := req.ProviderData.(ProviderData); ok {
		*out = providerData

		return diags
	}

	diags.AddError("Invalid provider data", "")

	return diags
}
