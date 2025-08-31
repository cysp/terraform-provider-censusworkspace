package provider

func (c *BigQueryDestinationCredentials) UpdateWithConnectionDetails(connectionDetails BigQueryDestinationConnectionDetails) {
	c.ProjectID = connectionDetails.ProjectID
	c.Location = connectionDetails.Location
}
