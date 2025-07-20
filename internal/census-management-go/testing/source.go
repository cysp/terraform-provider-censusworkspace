package testing

import (
	cm "github.com/cysp/terraform-provider-censusworkspace/internal/census-management-go"
)

func NewSourceFromCreateSourceData(ID int64, data cm.CreateSourceData) cm.SourceData {
	source := cm.SourceData{
		ID: ID,
	}

	UpdateSourceWithCreateSourceData(&source, data)

	return source
}

func UpdateSourceWithCreateSourceData(source *cm.SourceData, data cm.CreateSourceData) {
	source.Name = data.Type
	source.Type = data.Type

	if label, labelOk := data.Label.Get(); labelOk {
		source.Label.SetTo(label)
	} else {
		source.Label.SetToNull()
	}

	source.ConnectionDetails = data.Credentials
}

func UpdateSourceWithUpdateSourceData(source *cm.SourceData, data cm.UpdateSourceData) {
	if label, labelOk := data.Label.Get(); labelOk {
		source.Label.SetTo(label)
	} else {
		source.Label.SetToNull()
	}

	if data.Credentials != nil {
		source.ConnectionDetails = data.Credentials
	}
}
