package provider

import (
	"context"
	"maps"
	"slices"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
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

func (c resolvedStringCredential) AddValueTo(values writeOnlyCredentialValues) {
	if c.UsedWriteOnly {
		values.Add(c.WriteOnlyPath, c.Value)
	}
}

func resolveStringCredentialInput(input stringCredentialInput) (resolvedStringCredential, diag.Diagnostics) {
	value, usedWriteOnly, diags := resolveStringCredential(input.Legacy, input.WriteOnly, input.LegacyPath, input.WriteOnlyPath, input.Required)

	return resolvedStringCredential{
		Value:         value,
		UsedWriteOnly: usedWriteOnly,
		WriteOnlyPath: input.WriteOnlyPath,
	}, diags
}

func brazeDestinationModelWithWriteOnlyCredentials(plan, config BrazeDestinationModel) (BrazeDestinationModel, writeOnlyCredentialValues, diag.Diagnostics) {
	var diags diag.Diagnostics

	values := writeOnlyCredentialValues{}

	model := plan
	credentials := plan.Credentials.Value()
	configCredentials := config.Credentials.Value()

	credentials.InstanceURL = configCredentials.InstanceURL

	apiKey, apiKeyDiags := resolveStringCredentialInput(stringCredentialInput{
		Legacy:        configCredentials.APIKey,
		WriteOnly:     configCredentials.APIKeyWO,
		LegacyPath:    path.Root("credentials").AtName("api_key"),
		WriteOnlyPath: path.Root("credentials").AtName("api_key_wo"),
		Required:      true,
	})
	diags.Append(apiKeyDiags...)
	apiKey.AddValueTo(values)

	if apiKey.UsedWriteOnly {
		credentials.APIKeyWO = types.StringNull()
	}

	credentials.APIKey = apiKey.Value

	clientKeyValue := configCredentials.ClientKey
	if !writeOnlyStringConfigured(clientKeyValue) && writeOnlyStringConfigured(credentials.ClientKey) {
		clientKeyValue = credentials.ClientKey
	}

	clientKey, clientKeyDiags := resolveStringCredentialInput(stringCredentialInput{
		Legacy:        clientKeyValue,
		WriteOnly:     configCredentials.ClientKeyWO,
		LegacyPath:    path.Root("credentials").AtName("client_key"),
		WriteOnlyPath: path.Root("credentials").AtName("client_key_wo"),
	})
	diags.Append(clientKeyDiags...)
	clientKey.AddValueTo(values)

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

	if writeOnlyStringConfigured(configCredentials.APIKeyWO) {
		credentials.APIKey = types.StringNull()
	}

	if writeOnlyStringConfigured(configCredentials.ClientKeyWO) {
		credentials.ClientKey = types.StringNull()
	}

	credentials.APIKeyWO = types.StringNull()
	credentials.ClientKeyWO = types.StringNull()
	credentials.UpdateWithConnectionDetails(connectionDetails)

	model.Credentials = NewTypedObject(credentials)

	return model
}

func bigQueryDestinationModelWithWriteOnlyCredentials(plan, config BigQueryDestinationModel) (BigQueryDestinationModel, writeOnlyCredentialValues, diag.Diagnostics) {
	var diags diag.Diagnostics

	values := writeOnlyCredentialValues{}

	model := plan
	credentials := plan.Credentials.Value()
	configCredentials := config.Credentials.Value()

	credentials.ProjectID = configCredentials.ProjectID
	credentials.Location = configCredentials.Location

	serviceAccountKey, serviceAccountKeyDiags := resolveStringCredentialInput(stringCredentialInput{
		Legacy:        configCredentials.ServiceAccountKey,
		WriteOnly:     configCredentials.ServiceAccountKeyWO,
		LegacyPath:    path.Root("credentials").AtName("service_account_key"),
		WriteOnlyPath: path.Root("credentials").AtName("service_account_key_wo"),
	})
	diags.Append(serviceAccountKeyDiags...)
	serviceAccountKey.AddValueTo(values)

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

	if writeOnlyStringConfigured(configCredentials.ServiceAccountKeyWO) {
		credentials.ServiceAccountKey = types.StringNull()
	}

	credentials.ServiceAccountKeyWO = types.StringNull()
	credentials.UpdateWithConnectionDetails(connectionDetails)

	model.Credentials = NewTypedObject(credentials)

	return model
}

func bigQuerySourceModelWithWriteOnlyCredentials(plan, config BigQuerySourceModel) (BigQuerySourceModel, writeOnlyCredentialValues, diag.Diagnostics) {
	var diags diag.Diagnostics

	values := writeOnlyCredentialValues{}

	model := plan
	credentials := plan.Credentials.Value()
	configCredentials := config.Credentials.Value()

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
		privateKey.AddValueTo(values)

		serviceAccountKey.PrivateKey = privateKey.Value

		if privateKey.UsedWriteOnly {
			serviceAccountKey.PrivateKeyWO = types.StringNull()
		}

		credentials.ServiceAccountKey = NewTypedObject(serviceAccountKey)
	} else if configServiceAccountKeyOk && writeOnlyStringConfigured(configServiceAccountKey.PrivateKeyWO) {
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
		if configServiceAccountKeyOk && writeOnlyStringConfigured(configServiceAccountKey.PrivateKeyWO) {
			serviceAccountKey.PrivateKey = types.StringNull()
		}

		serviceAccountKey.PrivateKeyWO = types.StringNull()
		credentials.ServiceAccountKey = NewTypedObject(serviceAccountKey)
	}

	credentials.UpdateWithConnectionDetails(connectionDetails)
	model.Credentials = NewTypedObject(credentials)

	return model
}

func customAPIDestinationModelWithWriteOnlyCredentials(plan, config CustomAPIDestinationModel) (CustomAPIDestinationModel, writeOnlyCredentialValues, diag.Diagnostics) {
	var diags diag.Diagnostics

	values := writeOnlyCredentialValues{}

	model := plan
	credentials := config.Credentials.Value()
	configCredentials := config.Credentials.Value()

	if credentials.CustomHeaders.IsNull() || credentials.CustomHeaders.IsUnknown() {
		model.Credentials = NewTypedObject(credentials)

		return model, values, diags
	}

	headers := make(map[string]TypedObject[CustomAPIDestinationCustomHeader], len(credentials.CustomHeaders.Elements()))

	configHeaders := map[string]TypedObject[CustomAPIDestinationCustomHeader]{}
	if !configCredentials.CustomHeaders.IsNull() && !configCredentials.CustomHeaders.IsUnknown() {
		configHeaders = configCredentials.CustomHeaders.Elements()
	}

	for key, header := range credentials.CustomHeaders.Elements() {
		headerValue := header.Value()
		if configHeader, ok := configHeaders[key]; ok {
			configHeaderValue := configHeader.Value()
			writeOnlyPath := customAPIHeaderValueWOPath(key)
			value, valueDiags := resolveStringCredentialInput(stringCredentialInput{
				Legacy:        headerValue.Value,
				WriteOnly:     configHeaderValue.ValueWO,
				LegacyPath:    customAPIHeaderValuePath(key),
				WriteOnlyPath: writeOnlyPath,
				Required:      true,
			})
			diags.Append(valueDiags...)
			value.AddValueTo(values)

			headerValue.Value = value.Value
			if value.UsedWriteOnly {
				headerValue.ValueWO = types.StringNull()
			}
		}

		headerValue.ValueWO = types.StringNull()
		headers[key] = NewTypedObject(headerValue)
	}

	planHeaders := map[string]TypedObject[CustomAPIDestinationCustomHeader]{}
	if planCredentials := plan.Credentials.Value(); !planCredentials.CustomHeaders.IsNull() && !planCredentials.CustomHeaders.IsUnknown() {
		planHeaders = planCredentials.CustomHeaders.Elements()
	}

	for key, header := range planHeaders {
		if _, ok := headers[key]; ok {
			continue
		}

		headerValue := header.Value()
		headerValue.ValueWO = types.StringNull()
		headers[key] = NewTypedObject(headerValue)
	}

	credentials.CustomHeaders = NewTypedMap(headers)
	model.Credentials = NewTypedObject(credentials)

	return model, values, diags
}

func hydrateCustomAPIHeaderWriteOnlyValues(ctx context.Context, config tfsdk.Config, model *CustomAPIDestinationModel) diag.Diagnostics {
	var diags diag.Diagnostics

	credentials := model.Credentials.Value()
	if credentials.CustomHeaders.IsNull() || credentials.CustomHeaders.IsUnknown() {
		return diags
	}

	headers := maps.Clone(credentials.CustomHeaders.Elements())
	for key, header := range headers {
		var value types.String

		diags.Append(config.GetAttribute(ctx, customAPIHeaderValueWOPath(key), &value)...)

		if diags.HasError() {
			return diags
		}

		headerValue := header.Value()
		headerValue.ValueWO = value
		headers[key] = NewTypedObject(headerValue)
	}

	credentials.CustomHeaders = NewTypedMap(headers)
	model.Credentials = NewTypedObject(credentials)

	return diags
}

func customAPIHeaderValuePath(key string) path.Path {
	return path.Root("credentials").AtName("custom_headers").AtMapKey(key).AtName("value")
}

func customAPIHeaderValueWOPath(key string) path.Path {
	return path.Root("credentials").AtName("custom_headers").AtMapKey(key).AtName("value_wo")
}

func sanitizedCustomAPIDestinationCredentials(model, plan, config CustomAPIDestinationModel, connectionDetails CustomAPIDestinationConnectionDetails, writeOnlyValues writeOnlyCredentialValues) CustomAPIDestinationModel {
	credentials := plan.Credentials.Value()
	configCredentials := config.Credentials.Value()

	credentials.APIVersion = connectionDetails.APIVersion
	credentials.WebhookURL = connectionDetails.WebhookURL

	if connectionDetails.CustomHeaders.IsNull() || connectionDetails.CustomHeaders.IsUnknown() {
		credentials.CustomHeaders = NewTypedMapNull[TypedObject[CustomAPIDestinationCustomHeader]]()
		model.Credentials = NewTypedObject(credentials)

		return model
	}

	headers := mergeCustomAPIHeaders(
		credentials.CustomHeaders,
		configCredentials.CustomHeaders,
		connectionDetails.CustomHeaders,
		writeOnlyValues,
	)
	credentials.CustomHeaders = NewTypedMap(headers)
	model.Credentials = NewTypedObject(credentials)
	model.ConnectionDetails = NewTypedObject(sanitizedCustomAPIConnectionDetails(connectionDetails, credentials.CustomHeaders, writeOnlyValues))

	return model
}

func mergeCustomAPIHeaders(
	planned TypedMap[TypedObject[CustomAPIDestinationCustomHeader]],
	configured TypedMap[TypedObject[CustomAPIDestinationCustomHeader]],
	returned TypedMap[TypedObject[CustomAPIDestinationConnectionDetailsCustomHeader]],
	writeOnlyValues writeOnlyCredentialValues,
) map[string]TypedObject[CustomAPIDestinationCustomHeader] {
	plannedHeaders := map[string]TypedObject[CustomAPIDestinationCustomHeader]{}
	if !planned.IsNull() && !planned.IsUnknown() {
		plannedHeaders = planned.Elements()
	}

	configuredHeaders := map[string]TypedObject[CustomAPIDestinationCustomHeader]{}
	if !configured.IsNull() && !configured.IsUnknown() {
		configuredHeaders = configured.Elements()
	}

	headers := make(map[string]TypedObject[CustomAPIDestinationCustomHeader])

	keys := slices.Sorted(maps.Keys(returned.Elements()))
	for _, key := range keys {
		returnedHeader := returned.Elements()[key].Value()
		header := CustomAPIDestinationCustomHeader{
			Value:    returnedHeader.Value,
			IsSecret: returnedHeader.IsSecret,
			ValueWO:  types.StringNull(),
		}

		plannedHeader, plannedHeaderOk := plannedHeaders[key]
		configuredHeader, configuredHeaderOk := configuredHeaders[key]

		_, headerUsesWriteOnlyValue := writeOnlyValues[customAPIHeaderValueWOPath(key).String()]

		switch {
		case headerUsesWriteOnlyValue:
			header.Value = types.StringNull()
		case configuredHeaderOk && writeOnlyStringConfigured(configuredHeader.Value().ValueWO):
			header.Value = types.StringNull()
		case plannedHeaderOk && plannedHeader.Value().Value.IsNull():
			header.Value = types.StringNull()
		case header.Value.IsNull() && plannedHeaderOk:
			header.Value = plannedHeader.Value().Value
		}

		headers[key] = NewTypedObject(header)
	}

	return headers
}

func sanitizedCustomAPIConnectionDetails(
	connectionDetails CustomAPIDestinationConnectionDetails,
	credentialsHeaders TypedMap[TypedObject[CustomAPIDestinationCustomHeader]],
	writeOnlyValues writeOnlyCredentialValues,
) CustomAPIDestinationConnectionDetails {
	if connectionDetails.CustomHeaders.IsNull() || connectionDetails.CustomHeaders.IsUnknown() || credentialsHeaders.IsNull() || credentialsHeaders.IsUnknown() {
		return connectionDetails
	}

	credentialHeaders := credentialsHeaders.Elements()
	connectionHeaders := maps.Clone(connectionDetails.CustomHeaders.Elements())

	for key, credentialHeader := range credentialHeaders {
		connectionHeader, ok := connectionHeaders[key]
		if !ok {
			continue
		}

		_, headerUsesWriteOnlyValue := writeOnlyValues[customAPIHeaderValueWOPath(key).String()]
		if !headerUsesWriteOnlyValue && !credentialHeader.Value().Value.IsNull() {
			continue
		}

		connectionHeaderValue := connectionHeader.Value()
		connectionHeaderValue.Value = types.StringNull()
		connectionHeaders[key] = NewTypedObject(connectionHeaderValue)
	}

	connectionDetails.CustomHeaders = NewTypedMap(connectionHeaders)

	return connectionDetails
}
