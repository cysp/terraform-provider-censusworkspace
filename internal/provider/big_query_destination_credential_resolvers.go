package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func bigQueryDestinationModelWithWriteOnlyCredentials(plan, config BigQueryDestinationModel) (BigQueryDestinationModel, writeOnlyCredentialValues, diag.Diagnostics) {
	var diags diag.Diagnostics

	values := writeOnlyCredentialValues{}

	model := plan
	credentials := plan.Credentials.Value()
	configCredentials := configuredObjectOrPlan(config.Credentials, plan.Credentials).Value()

	credentials.ProjectID = configCredentials.ProjectID
	credentials.Location = configCredentials.Location

	serviceAccountKey, serviceAccountKeyDiags := resolveStringCredentialInput(stringCredentialInput{
		Legacy:        configCredentials.ServiceAccountKey,
		WriteOnly:     configCredentials.ServiceAccountKeyWO,
		LegacyPath:    path.Root("credentials").AtName("service_account_key"),
		WriteOnlyPath: path.Root("credentials").AtName("service_account_key_wo"),
	})
	diags.Append(serviceAccountKeyDiags...)
	diags.Append(serviceAccountKey.AddValueTo(values)...)

	if serviceAccountKey.UsedWriteOnly {
		credentials.ServiceAccountKeyWO = types.StringNull()
	}

	credentials.ServiceAccountKey = serviceAccountKey.Value

	model.Credentials = NewTypedObject(credentials)

	return model, values, diags
}

func sanitizedBigQueryDestinationCredentials(model, plan, config BigQueryDestinationModel, connectionDetails BigQueryDestinationConnectionDetails) BigQueryDestinationModel {
	credentials := plan.Credentials.Value()
	configCredentials := config.Credentials.Value()

	if !configCredentials.ServiceAccountKeyWO.IsNull() {
		credentials.ServiceAccountKey = types.StringNull()
		connectionDetails.ServiceAccountKey = types.StringNull()
	}

	credentials.ServiceAccountKeyWO = types.StringNull()
	credentials.UpdateWithConnectionDetails(connectionDetails)

	model.Credentials = NewTypedObject(credentials)
	model.ConnectionDetails = NewTypedObject(connectionDetails)

	return model
}
