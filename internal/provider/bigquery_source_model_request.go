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
	m.Credentials.Encode(&enc)

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
	m.Credentials.Encode(&enc)

	body.Connection.Credentials = enc.Bytes()

	return body, nil
}

func (s *BigQuerySourceCredentials) Encode(e *jx.Encoder) {
	e.Obj(func(e *jx.Encoder) {
		e.Field("project_id", func(e *jx.Encoder) {
			e.Str(s.ProjectID.ValueString())
		})
		e.Field("location", func(e *jx.Encoder) {
			e.Str(s.Location.ValueString())
		})
	})
}
