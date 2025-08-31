package provider

import (
	"context"

	cm "github.com/cysp/terraform-provider-censusworkspace/internal/census-management-go"
	"github.com/go-faster/jx"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func (m *BigQuerySourceModel) ToCreateSourceData(_ context.Context) (cm.CreateSourceBody, diag.Diagnostics) {
	body := cm.CreateSourceBody{
		Connection: cm.CreateSourceBodyConnection{
			Type: BigQuerySourceType,
		},
	}

	label := m.Label.ValueString()
	if label != "" {
		body.Connection.Label.SetTo(label)
	} else {
		body.Connection.Label.SetToNull()
	}

	enc := jx.Encoder{}
	m.Credentials.Value().Encode(&enc)

	body.Connection.Credentials = enc.Bytes()

	return body, nil
}

func (m *BigQuerySourceModel) ToUpdateSourceData(_ context.Context) (cm.UpdateSourceBody, diag.Diagnostics) {
	body := cm.UpdateSourceBody{}

	label := m.Label.ValueString()
	if label != "" {
		body.Connection.Label.SetTo(label)
	} else {
		body.Connection.Label.SetToNull()
	}

	enc := jx.Encoder{}
	m.Credentials.Value().Encode(&enc)

	body.Connection.Credentials = enc.Bytes()

	return body, nil
}

func (c BigQuerySourceCredentials) Encode(enc *jx.Encoder) {
	enc.Obj(func(enc *jx.Encoder) {
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
