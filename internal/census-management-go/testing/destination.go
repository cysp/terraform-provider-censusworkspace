package testing

import (
	cm "github.com/cysp/terraform-provider-censusworkspace/internal/census-management-go"
)

func NewDestinationFromCreateDestinationBody(id int64, body cm.CreateDestinationBody) cm.DestinationData {
	destination := cm.DestinationData{
		ID: id,
	}

	UpdateDestinationWithCreateDestinationBody(&destination, body)

	return destination
}

func UpdateDestinationWithCreateDestinationBody(destination *cm.DestinationData, body cm.CreateDestinationBody) {
	connection := body.ServiceConnection

	destination.Name = connection.Name
	destination.Type = connection.Type
	destination.ConnectionDetails = connection.Credentials
}

func UpdateDestinationWithUpdateDestinationBody(destination *cm.DestinationData, body cm.UpdateDestinationBody) {
	connection := body.ServiceConnection

	if name, nameOk := connection.Name.Get(); nameOk {
		destination.Name = name
	}

	if connection.Credentials != nil {
		destination.ConnectionDetails = connection.Credentials
	}
}
