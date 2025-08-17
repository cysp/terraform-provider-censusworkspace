package testing

import (
	"time"

	cm "github.com/cysp/terraform-provider-censusworkspace/internal/census-management-go"
)

func NewSourceModelFromCreateSourceModelBody(ID int64, body cm.CreateSourceModelBody) cm.SourceModelData {
	model := cm.SourceModelData{
		Type: cm.SourceModelDataTypeModel,
		ID:   ID,
	}

	UpdateSourceModelWithCreateSourceModelBody(&model, body)

	return model
}

func UpdateSourceModelWithCreateSourceModelBody(model *cm.SourceModelData, body cm.CreateSourceModelBody) {
	model.Name = body.Name
	model.Query = body.Query

	model.Description = body.Description

	model.CreatedAt = time.Now()
	model.UpdatedAt.SetTo(time.Now())
}

func UpdateSourceModelWithUpdateSourceModelBody(model *cm.SourceModelData, body cm.UpdateSourceModelBody) {
	if name, nameOk := body.Name.Get(); nameOk {
		model.Name = name
	}

	if query, queryOk := body.Query.Get(); queryOk {
		model.Query = query
	}

	if description, descriptionOk := body.Description.Get(); descriptionOk {
		model.Description.SetTo(description)
	}

	model.UpdatedAt.SetTo(time.Now())
}
