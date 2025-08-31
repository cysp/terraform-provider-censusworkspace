package provider

func (c *BrazeDestinationCredentials) UpdateWithConnectionDetails(connectionDetails BrazeDestinationConnectionDetails) {
	c.InstanceURL = connectionDetails.InstanceURL
}
