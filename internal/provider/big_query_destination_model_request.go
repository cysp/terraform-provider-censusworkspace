package provider

import (
	"context"

	cm "github.com/cysp/terraform-provider-censusworkspace/internal/census-management-go"
	"github.com/go-faster/jx"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
)

func (m *BigQueryDestinationModel) ToCreateDestinationData(_ context.Context) (cm.CreateDestinationBody, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	body := cm.CreateDestinationBody{
		ServiceConnection: cm.CreateDestinationBodyServiceConnection{
			Name: m.Name.ValueString(),
			Type: BigQueryDestinationType,
		},
	}

	enc := jx.Encoder{}
	if !m.Credentials.IsNull() && !m.Credentials.IsUnknown() {
		credentialsEncodeFailed := m.Credentials.Value().Encode(&enc)
		if credentialsEncodeFailed {
			diags.AddAttributeError(path.Root("credentials"), "Failed to encode value", "")
		}

		body.ServiceConnection.Credentials = enc.Bytes()
	}

	return body, diags
}

func (m *BigQueryDestinationModel) ToUpdateDestinationData(_ context.Context) (cm.UpdateDestinationBody, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	body := cm.UpdateDestinationBody{
		ServiceConnection: cm.UpdateDestinationBodyServiceConnection{
			Name: cm.NewOptString(m.Name.ValueString()),
		},
	}

	enc := jx.Encoder{}
	if !m.Credentials.IsNull() && !m.Credentials.IsUnknown() {
		credentialsEncodeFailed := m.Credentials.Value().Encode(&enc)
		if credentialsEncodeFailed {
			diags.AddAttributeError(path.Root("credentials"), "Failed to encode value", "")
		}

		body.ServiceConnection.Credentials = enc.Bytes()
	}

	return body, diags
}

func (c BigQueryDestinationCredentials) Encode(enc *jx.Encoder) bool {
	return enc.Obj(func(enc *jx.Encoder) {
		enc.Field("project_id", func(e *jx.Encoder) {
			e.Str(c.ProjectID.ValueString())
		})

		enc.Field("location", func(e *jx.Encoder) {
			e.Str(c.Location.ValueString())
		})

		serviceAccountKey := c.ServiceAccountKey.ValueString()
		if serviceAccountKey != "" {
			enc.Field("service_account_key", func(e *jx.Encoder) {
				e.Str(serviceAccountKey)
			})
		}
	})
}
