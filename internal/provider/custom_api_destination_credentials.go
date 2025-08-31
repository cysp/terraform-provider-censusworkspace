package provider

func (c *CustomAPIDestinationCredentials) UpdateWithConnectionDetails(connectionDetails CustomAPIDestinationConnectionDetails) {
	c.APIVersion = connectionDetails.APIVersion
	c.WebhookURL = connectionDetails.WebhookURL

	if connectionDetails.CustomHeaders.IsNull() || connectionDetails.CustomHeaders.IsUnknown() {
		c.CustomHeaders = connectionDetails.CustomHeaders
	} else {
		existingCustomHeaders := c.CustomHeaders.Elements()

		customHeaders := make(map[string]TypedObject[CustomAPIDestinationCustomHeader])

		for key, value := range connectionDetails.CustomHeaders.Elements() {
			existingHeader, existingHeaderOk := existingCustomHeaders[key]

			customHeader := value.Value()

			if customHeader.Value.IsNull() && existingHeaderOk {
				customHeader.Value = existingHeader.Value().Value
			}

			customHeaders[key] = NewTypedObject(customHeader)
		}

		c.CustomHeaders = NewTypedMap(customHeaders)
	}
}
