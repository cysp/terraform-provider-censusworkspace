package testing

import (
	cm "github.com/cysp/terraform-provider-censusworkspace/internal/census-management-go"
)

func NewSourceFromCreateSourceBody(id int64, body cm.CreateSourceBody) cm.SourceData {
	source := cm.SourceData{
		ID: id,
	}

	UpdateSourceWithCreateSourceBody(&source, body)

	return source
}

func UpdateSourceWithCreateSourceBody(source *cm.SourceData, body cm.CreateSourceBody) {
	connection := body.Connection

	source.Name = connection.Type
	source.Type = connection.Type

	if label, labelOk := connection.Label.Get(); labelOk {
		source.Label.SetTo(label)
	} else {
		source.Label.SetToNull()
	}

	source.ConnectionDetails = connection.Credentials
}

func UpdateSourceWithUpdateSourceBody(source *cm.SourceData, body cm.UpdateSourceBody) {
	connection := body.Connection

	if label, labelOk := connection.Label.Get(); labelOk {
		source.Label.SetTo(label)
	} else {
		source.Label.SetToNull()
	}

	if connection.Credentials != nil {
		source.ConnectionDetails = connection.Credentials
	}
}
