package provider

import (
	"context"

	cm "github.com/cysp/terraform-provider-censusworkspace/internal/census-management-go"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func (m *SyncModel) ToCreateOrUpdateSyncData(_ context.Context) (cm.CreateOrUpdateSyncBody, diag.Diagnostics) {
	var diags diag.Diagnostics

	createOrUpdateSyncBody := cm.CreateOrUpdateSyncBody{}

	if !m.Label.IsNull() {
		createOrUpdateSyncBody.Label.SetTo(m.Label.ValueString())
	} else {
		createOrUpdateSyncBody.Label.SetToNull()
	}

	return createOrUpdateSyncBody, diags
}
