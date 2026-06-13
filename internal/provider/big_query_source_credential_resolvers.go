package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func bigQuerySourceModelWithWriteOnlyCredentials(plan, config BigQuerySourceModel) (BigQuerySourceModel, writeOnlyCredentialValues, diag.Diagnostics) {
	var diags diag.Diagnostics

	values := writeOnlyCredentialValues{}

	model := plan
	credentials := plan.Credentials.Value()
	configCredentials := configuredObjectOrPlan(config.Credentials, plan.Credentials).Value()

	credentials.ProjectID = configCredentials.ProjectID
	credentials.Location = configCredentials.Location
	credentials.ServiceAccountKey = configCredentials.ServiceAccountKey

	serviceAccountKey, serviceAccountKeyOk := credentials.ServiceAccountKey.GetValue()
	configServiceAccountKey, configServiceAccountKeyOk := configCredentials.ServiceAccountKey.GetValue()

	if serviceAccountKeyOk {
		privateKey, privateKeyDiags := resolveStringCredentialInput(stringCredentialInput{
			Legacy:        configServiceAccountKey.PrivateKey,
			WriteOnly:     configServiceAccountKey.PrivateKeyWO,
			LegacyPath:    path.Root("credentials").AtName("service_account_key").AtName("private_key"),
			WriteOnlyPath: path.Root("credentials").AtName("service_account_key").AtName("private_key_wo"),
			Required:      true,
		})
		diags.Append(privateKeyDiags...)
		diags.Append(privateKey.AddValueTo(values)...)

		serviceAccountKey.PrivateKey = privateKey.Value

		if privateKey.UsedWriteOnly {
			serviceAccountKey.PrivateKeyWO = types.StringNull()
		}

		credentials.ServiceAccountKey = NewTypedObject(serviceAccountKey)
	} else if configServiceAccountKeyOk && !configServiceAccountKey.PrivateKeyWO.IsNull() {
		diags.AddError(
			"Missing credential argument",
			path.Root("credentials").AtName("service_account_key").AtName("private_key_wo").String()+" requires "+
				path.Root("credentials").AtName("service_account_key").String()+" to be configured.",
		)
	}

	model.Credentials = NewTypedObject(credentials)

	return model, values, diags
}

func sanitizedBigQuerySourceCredentials(model, plan, config BigQuerySourceModel, connectionDetails BigQuerySourceConnectionDetails) BigQuerySourceModel {
	credentials := plan.Credentials.Value()
	configCredentials := config.Credentials.Value()

	serviceAccountKey, serviceAccountKeyOk := credentials.ServiceAccountKey.GetValue()
	configServiceAccountKey, configServiceAccountKeyOk := configCredentials.ServiceAccountKey.GetValue()

	if serviceAccountKeyOk {
		if configServiceAccountKeyOk && !configServiceAccountKey.PrivateKeyWO.IsNull() {
			serviceAccountKey.PrivateKey = types.StringNull()
		}

		serviceAccountKey.PrivateKeyWO = types.StringNull()
		credentials.ServiceAccountKey = NewTypedObject(serviceAccountKey)
	}

	credentials.UpdateWithConnectionDetails(connectionDetails)
	model.Credentials = NewTypedObject(credentials)

	return model
}
