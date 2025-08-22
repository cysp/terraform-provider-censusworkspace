package provider

import (
	"context"

	cm "github.com/cysp/terraform-provider-censusworkspace/internal/census-management-go"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func (m *DestinationModel) ToCreateDestinationData(_ context.Context) (cm.CreateDestinationBody, diag.Diagnostics) {
	body := cm.CreateDestinationBody{
		ServiceConnection: cm.CreateDestinationBodyServiceConnection{
			Name: m.Name.ValueString(),
			Type: m.Type.ValueString(),
		},
	}

	if credentials := m.Credentials.ValueStringPointer(); credentials != nil {
		body.ServiceConnection.Credentials = []byte(*credentials)
	}

	return body, nil
}

func (m *DestinationModel) ToUpdateDestinationData(_ context.Context) (cm.UpdateDestinationBody, diag.Diagnostics) {
	body := cm.UpdateDestinationBody{
		ServiceConnection: cm.UpdateDestinationBodyServiceConnection{
			Name: cm.NewOptString(m.Name.ValueString()),
		},
	}

	if credentials := m.Credentials.ValueStringPointer(); credentials != nil {
		body.ServiceConnection.Credentials = []byte(*credentials)
	}

	return body, nil
}
