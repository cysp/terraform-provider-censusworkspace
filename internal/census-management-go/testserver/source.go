package testserver

import (
	cm "github.com/cysp/terraform-provider-censusworkspace/internal/census-management-go"
)

func NewSourceFromCreateSourceData(ID int64, data cm.CreateSourceData) cm.SourceData {
	source := cm.SourceData{
		ID: ID,
	}

	UpdateSourceFromCreateSourceData(&source, data)

	return source
}

func UpdateSourceFromCreateSourceData(source *cm.SourceData, data cm.CreateSourceData) {
	source.Name = data.Type
	source.Type = data.Type

	if label, labelOk := data.Label.Get(); labelOk {
		source.Label.SetTo(label)
	} else {
		source.Label.SetToNull()
	}

	source.ConnectionDetails = data.Credentials
}

func UpdateSourceFromUpdateSourceData(source *cm.SourceData, data cm.UpdateSourceData) {
	if label, labelOk := data.Label.Get(); labelOk {
		source.Label.SetTo(label)
	} else {
		source.Label.SetToNull()
	}

	if data.Credentials != nil {
		source.ConnectionDetails = data.Credentials
	}
}
