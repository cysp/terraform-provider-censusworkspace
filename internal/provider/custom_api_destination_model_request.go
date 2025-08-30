package provider

import (
	"context"

	cm "github.com/cysp/terraform-provider-censusworkspace/internal/census-management-go"
	"github.com/go-faster/jx"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
)

func (m *CustomAPIDestinationModel) ToCreateDestinationData(_ context.Context) (cm.CreateDestinationBody, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	body := cm.CreateDestinationBody{
		ServiceConnection: cm.CreateDestinationBodyServiceConnection{
			Name: m.Name.ValueString(),
			Type: CustomAPIDestinationType,
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

func (m *CustomAPIDestinationModel) ToUpdateDestinationData(_ context.Context) (cm.UpdateDestinationBody, diag.Diagnostics) {
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

func (c CustomAPIDestinationCredentials) Encode(enc *jx.Encoder) bool {
	return enc.Obj(func(enc *jx.Encoder) {
		if !c.APIVersion.IsNull() && !c.APIVersion.IsUnknown() {
			enc.Field("api_version", func(enc *jx.Encoder) {
				enc.Int64(c.APIVersion.ValueInt64())
			})
		}

		if !c.WebhookURL.IsNull() && !c.WebhookURL.IsUnknown() {
			enc.Field("webhook_url", func(enc *jx.Encoder) {
				enc.Str(c.WebhookURL.ValueString())
			})
		}

		if !c.CustomHeaders.IsNull() && !c.CustomHeaders.IsUnknown() {
			enc.Field("custom_headers", func(enc *jx.Encoder) {
				enc.Obj(func(enc *jx.Encoder) {
					for k, v := range c.CustomHeaders.Elements() {
						enc.Field(k, func(enc *jx.Encoder) {
							v.Value().Encode(enc)
						})
					}
				})
			})
		}
	})
}

func (c CustomAPIDestinationCustomHeader) Encode(enc *jx.Encoder) bool {
	return enc.Obj(func(enc *jx.Encoder) {
		if !c.Value.IsNull() && !c.Value.IsUnknown() {
			enc.Field("value", func(enc *jx.Encoder) {
				enc.Str(c.Value.ValueString())
			})
		}

		if !c.IsSecret.IsNull() && !c.IsSecret.IsUnknown() {
			enc.Field("is_secret", func(enc *jx.Encoder) {
				enc.Bool(c.IsSecret.ValueBool())
			})
		}
	})
}
