package testing

import (
	cm "github.com/cysp/terraform-provider-censusworkspace/internal/census-management-go"
)

func NewSyncFromCreateSyncBody(id int64, body cm.CreateSyncBody) cm.SyncData {
	sync := cm.SyncData{
		ID: id,
	}

	UpdateSyncWithCreateSyncBody(&sync, body)

	return sync
}

func UpdateSyncWithCreateSyncBody(sync *cm.SyncData, body cm.CreateSyncBody) {
	if label, labelOk := body.Label.Get(); labelOk {
		sync.Label.SetTo(label)
	} else {
		sync.Label.SetToNull()
	}

	sync.Operation = body.Operation

	sync.DestinationAttributes = cm.SyncDataDestinationAttributes(body.DestinationAttributes)
}

func UpdateSyncWithUpdateSyncBody(sync *cm.SyncData, body cm.UpdateSyncBody) {
	if label, labelOk := body.Label.Get(); labelOk {
		sync.Label.SetTo(label)
	} else {
		sync.Label.SetToNull()
	}

	sync.Operation = body.Operation

	sync.DestinationAttributes = cm.SyncDataDestinationAttributes(body.DestinationAttributes)
}
