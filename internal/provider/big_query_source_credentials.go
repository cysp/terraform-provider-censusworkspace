package provider

func (c *BigQuerySourceCredentials) UpdateWithConnectionDetails(connectionDetails BigQuerySourceConnectionDetails) {
	c.ProjectID = connectionDetails.ProjectID
	c.Location = connectionDetails.Location
}
