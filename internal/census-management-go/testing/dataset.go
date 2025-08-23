package testing

import (
	"errors"
	"fmt"
	"time"

	cm "github.com/cysp/terraform-provider-censusworkspace/internal/census-management-go"
)

func NewDatasetFromCreateDatasetBody(id int64, body cm.CreateDatasetBody) cm.DatasetData {
	dataset := cm.DatasetData{}

	//nolint:gocritic
	switch body.Type {
	case cm.CreateSQLDatasetBodyCreateDatasetBody:
		sqlBody, ok := body.GetCreateSQLDatasetBody()
		if ok {
			dataset.SetSQLDatasetData(cm.SQLDatasetData{
				ID:          id,
				Type:        cm.SQLDatasetDataTypeSQL,
				Name:        sqlBody.Name,
				SourceID:    sqlBody.SourceID,
				Query:       sqlBody.Query,
				Description: sqlBody.Description,
			})
		}
	}

	UpdateDatasetWithCreateDatasetBody(&dataset, body)

	return dataset
}

func UpdateDatasetWithCreateDatasetBody(dataset *cm.DatasetData, body cm.CreateDatasetBody) {
	//nolint:gocritic
	switch dataset.Type {
	case cm.SQLDatasetDataDatasetData:
		sql := &dataset.SQLDatasetData

		body, bodyOk := body.GetCreateSQLDatasetBody()
		if !bodyOk {
			return
		}

		sql.Name = body.Name
		sql.Query = body.Query
		sql.SourceID = body.SourceID
		sql.Description = body.Description

		now := time.Now()
		sql.CreatedAt = now
		sql.UpdatedAt = now
	}
}

func UpdateDatasetWithUpdateDatasetBody(dataset *cm.DatasetData, body cm.UpdateDatasetBody) error {
	switch dataset.Type {
	case cm.SQLDatasetDataDatasetData:
		dataset := &dataset.SQLDatasetData

		body, bodyOk := body.GetUpdateSQLDatasetBody()
		if !bodyOk {
			return errDatasetInvalidRequestBody
		}

		if name, nameOk := body.Name.Get(); nameOk {
			dataset.Name = name
		}

		if query, queryOk := body.Query.Get(); queryOk {
			dataset.Query = query
		}

		if description, descriptionOk := body.Description.Get(); descriptionOk {
			dataset.Description.SetTo(description)
		}

		dataset.UpdatedAt = time.Now()

		return nil
	default:
		return fmt.Errorf("%w: %v", errDatasetTypeUnknown, dataset.Type)
	}
}

var errDatasetInvalidRequestBody = errors.New("invalid request body")

var errDatasetTypeUnknown = errors.New("unknown dataset type")
