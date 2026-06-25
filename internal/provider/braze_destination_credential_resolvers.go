package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func brazeDestinationModelWithWriteOnlyCredentials(plan, config BrazeDestinationModel) (BrazeDestinationModel, writeOnlyCredentialValues, diag.Diagnostics) {
	var diags diag.Diagnostics

	values := writeOnlyCredentialValues{}

	model := plan
	credentials := plan.Credentials.Value()
	configCredentials := configuredObjectOrPlan(config.Credentials, plan.Credentials).Value()

	credentials.InstanceURL = configCredentials.InstanceURL

	apiKeyValue := configCredentials.APIKey
	if config.Credentials.IsUnknown() && stringKnown(credentials.APIKey) {
		apiKeyValue = credentials.APIKey
	}

	apiKey, apiKeyDiags := resolveStringCredentialInput(stringCredentialInput{
		Legacy:        apiKeyValue,
		WriteOnly:     configCredentials.APIKeyWO,
		LegacyPath:    path.Root("credentials").AtName("api_key"),
		WriteOnlyPath: path.Root("credentials").AtName("api_key_wo"),
		Required:      true,
	})
	diags.Append(apiKeyDiags...)
	diags.Append(apiKey.AddValueTo(values)...)

	if apiKey.UsedWriteOnly {
		credentials.APIKeyWO = types.StringNull()
	}

	credentials.APIKey = apiKey.Value

	clientKeyValue := configCredentials.ClientKey
	if !stringKnown(clientKeyValue) && stringKnown(credentials.ClientKey) {
		clientKeyValue = credentials.ClientKey
	}

	clientKey, clientKeyDiags := resolveStringCredentialInput(stringCredentialInput{
		Legacy:        clientKeyValue,
		WriteOnly:     configCredentials.ClientKeyWO,
		LegacyPath:    path.Root("credentials").AtName("client_key"),
		WriteOnlyPath: path.Root("credentials").AtName("client_key_wo"),
	})
	diags.Append(clientKeyDiags...)
	diags.Append(clientKey.AddValueTo(values)...)

	if clientKey.UsedWriteOnly {
		credentials.ClientKeyWO = types.StringNull()
	}

	credentials.ClientKey = clientKey.Value

	model.Credentials = NewTypedObject(credentials)

	return model, values, diags
}

func sanitizedBrazeDestinationCredentials(model, plan, config BrazeDestinationModel, connectionDetails BrazeDestinationConnectionDetails) BrazeDestinationModel {
	credentials := plan.Credentials.Value()
	configCredentials := config.Credentials.Value()

	if !configCredentials.APIKeyWO.IsNull() {
		credentials.APIKey = types.StringNull()
	}

	if !configCredentials.ClientKeyWO.IsNull() {
		credentials.ClientKey = types.StringNull()
	}

	credentials.APIKeyWO = types.StringNull()
	credentials.ClientKeyWO = types.StringNull()
	credentials.UpdateWithConnectionDetails(connectionDetails)

	model.Credentials = NewTypedObject(credentials)

	return model
}
