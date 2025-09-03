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

	syncEngine := m.SyncEngine.ValueString()
	if syncEngine != "" {
		body.Connection.SyncEngine.SetTo(syncEngine)
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

		serviceAccountKey, serviceAccountKeyOk := c.ServiceAccountKey.GetValue()
		if serviceAccountKeyOk {
			enc.Field("service_account_key", func(e *jx.Encoder) {
				serviceAccountKey.Encode(e)
			})
		}
	})
}

func (c BigQuerySourceCredentialsServiceAccountKey) Encode(enc *jx.Encoder) {
	enc.Obj(func(enc *jx.Encoder) {
		enc.Field("project_id", func(e *jx.Encoder) {
			e.Str(c.ProjectID.ValueString())
		})

		enc.Field("private_key_id", func(e *jx.Encoder) {
			e.Str(c.PrivateKeyID.ValueString())
		})

		enc.Field("private_key", func(e *jx.Encoder) {
			e.Str(c.PrivateKey.ValueString())
		})

		enc.Field("client_email", func(e *jx.Encoder) {
			e.Str(c.ClientEmail.ValueString())
		})

		enc.Field("client_id", func(e *jx.Encoder) {
			e.Str(c.ClientID.ValueString())
		})
	})
}
