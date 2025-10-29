package testing

import (
	cm "github.com/cysp/terraform-provider-censusworkspace/internal/census-management-go"
)

func NewSyncFromCreateSyncBody(id int64, body cm.CreateOrUpdateSyncBody) cm.SyncData {
	sync := cm.SyncData{
		ID: id,
	}

	UpdateSyncWithCreateOrUpdateSyncBody(&sync, body)

	return sync
}

func UpdateSyncWithCreateOrUpdateSyncBody(sync *cm.SyncData, body cm.CreateOrUpdateSyncBody) {
	if label, labelOk := body.Label.Get(); labelOk {
		sync.Label.SetTo(label)
	} else {
		sync.Label.SetToNull()
	}

	sync.Operation = body.Operation

	sync.DestinationAttributes = cm.SyncDestinationAttributes(body.DestinationAttributes)
}
