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

func customAPIDestinationModelWithWriteOnlyCredentials(plan, config CustomAPIDestinationModel) (CustomAPIDestinationModel, writeOnlyCredentialValues, diag.Diagnostics) {
	var diags diag.Diagnostics

	values := writeOnlyCredentialValues{}

	model := plan
	credentials := plan.Credentials.Value()
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
			diags.Append(value.AddValueTo(values)...)

			headerValue.Value = value.Value
			if value.UsedWriteOnly {
				headerValue.ValueWO = types.StringNull()
			}
		}

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

		headerUsesWriteOnlyValue, pathDiags := writeOnlyValues.Configured(customAPIHeaderValueWOPath(key))
		if pathDiags.HasError() {
			continue
		}

		switch {
		case headerUsesWriteOnlyValue:
			header.Value = types.StringNull()
		case configuredHeaderOk && !configuredHeader.Value().ValueWO.IsNull():
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

		headerUsesWriteOnlyValue, pathDiags := writeOnlyValues.Configured(customAPIHeaderValueWOPath(key))
		if pathDiags.HasError() {
			continue
		}

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
