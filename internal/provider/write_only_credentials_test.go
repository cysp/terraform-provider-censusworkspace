package provider //nolint:testpackage

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBrazeDestinationWriteOnlyCredentialResolver(t *testing.T) {
	t.Parallel()

	plan := BrazeDestinationModel{
		Credentials: NewTypedObject(BrazeDestinationCredentials{
			InstanceURL: types.StringValue("instance-url"),
			APIKey:      types.StringNull(),
			ClientKey:   types.StringValue("client-key"),
		}),
	}
	config := BrazeDestinationModel{
		Credentials: NewTypedObject(BrazeDestinationCredentials{
			APIKeyWO: types.StringValue("api-key"),
		}),
	}

	requestModel, values, diags := brazeDestinationModelWithWriteOnlyCredentials(plan, config)
	require.False(t, diags.HasError(), diags)
	assert.Contains(t, values, path.Root("credentials").AtName("api_key_wo").String())
	assert.Equal(t, "api-key", requestModel.Credentials.Value().APIKey.ValueString())
	assert.Equal(t, "client-key", requestModel.Credentials.Value().ClientKey.ValueString())

	stateModel := sanitizedBrazeDestinationCredentials(plan, plan, config, BrazeDestinationConnectionDetails{
		InstanceURL: types.StringValue("normalized-instance-url"),
	})
	assert.True(t, stateModel.Credentials.Value().APIKey.IsNull())
	assert.Equal(t, "client-key", stateModel.Credentials.Value().ClientKey.ValueString())
	assert.Equal(t, "normalized-instance-url", stateModel.Credentials.Value().InstanceURL.ValueString())
}

func TestCustomAPIDestinationWriteOnlyHeaderResolverKeepsIsSecretOrthogonal(t *testing.T) {
	t.Parallel()

	plan := CustomAPIDestinationModel{
		Credentials: NewTypedObject(CustomAPIDestinationCredentials{
			CustomHeaders: NewTypedMap(map[string]TypedObject[CustomAPIDestinationCustomHeader]{
				"x-client-id": NewTypedObject(CustomAPIDestinationCustomHeader{
					Value:    types.StringValue("client-id"),
					IsSecret: types.BoolValue(false),
				}),
				"x-client-secret": NewTypedObject(CustomAPIDestinationCustomHeader{
					Value:    types.StringNull(),
					IsSecret: types.BoolValue(true),
				}),
			}),
		}),
	}
	config := CustomAPIDestinationModel{
		Credentials: NewTypedObject(CustomAPIDestinationCredentials{
			CustomHeaders: NewTypedMap(map[string]TypedObject[CustomAPIDestinationCustomHeader]{
				"x-client-secret": NewTypedObject(CustomAPIDestinationCustomHeader{
					ValueWO:  types.StringValue("secret"),
					IsSecret: types.BoolValue(true),
				}),
			}),
		}),
	}

	requestModel, values, diags := customAPIDestinationModelWithWriteOnlyCredentials(plan, config)
	require.False(t, diags.HasError(), diags)
	assert.Contains(t, values, path.Root("credentials").AtName("custom_headers").AtMapKey("x-client-secret").AtName("value_wo").String())

	headers := requestModel.Credentials.Value().CustomHeaders.Elements()
	assert.Equal(t, "client-id", headers["x-client-id"].Value().Value.ValueString())
	assert.False(t, headers["x-client-id"].Value().IsSecret.ValueBool())
	assert.Equal(t, "secret", headers["x-client-secret"].Value().Value.ValueString())
	assert.True(t, headers["x-client-secret"].Value().IsSecret.ValueBool())

	stateModel := sanitizedCustomAPIDestinationCredentials(
		plan,
		plan,
		config,
		CustomAPIDestinationConnectionDetails{
			CustomHeaders: NewTypedMap(map[string]TypedObject[CustomAPIDestinationConnectionDetailsCustomHeader]{
				"x-client-id": NewTypedObject(CustomAPIDestinationConnectionDetailsCustomHeader{
					Value:    types.StringValue("client-id"),
					IsSecret: types.BoolValue(false),
				}),
				"x-client-secret": NewTypedObject(CustomAPIDestinationConnectionDetailsCustomHeader{
					Value:    types.StringNull(),
					IsSecret: types.BoolValue(true),
				}),
			}),
		},
		values,
	)
	stateHeaders := stateModel.Credentials.Value().CustomHeaders.Elements()
	assert.True(t, stateHeaders["x-client-secret"].Value().Value.IsNull())
	assert.True(t, stateHeaders["x-client-secret"].Value().IsSecret.ValueBool())
}

func TestWriteOnlyCredentialVerifierIncludesCanonicalPath(t *testing.T) {
	t.Parallel()

	value := types.StringValue("shared-secret")
	apiKeyPath := path.Root("credentials").AtName("api_key_wo")
	clientKeyPath := path.Root("credentials").AtName("client_key_wo")
	apiKeyVerifier, err := writeOnlyCredentialVerifier(apiKeyPath.String(), value)
	require.NoError(t, err)

	assert.True(t, writeOnlyCredentialVerifierMatches(apiKeyPath.String(), value, apiKeyVerifier))
	assert.False(t, writeOnlyCredentialVerifierMatches(clientKeyPath.String(), value, apiKeyVerifier))
}

func TestCustomAPIHeaderWriteOnlyPathEscapesMapKeys(t *testing.T) {
	t.Parallel()

	key := `x-"secret"\key`
	expectedPath := path.Root("credentials").AtName("custom_headers").AtMapKey(key).AtName("value_wo")

	plan := CustomAPIDestinationModel{
		Credentials: NewTypedObject(CustomAPIDestinationCredentials{
			CustomHeaders: NewTypedMap(map[string]TypedObject[CustomAPIDestinationCustomHeader]{
				key: NewTypedObject(CustomAPIDestinationCustomHeader{
					Value:    types.StringNull(),
					IsSecret: types.BoolValue(true),
				}),
			}),
		}),
	}
	config := CustomAPIDestinationModel{
		Credentials: NewTypedObject(CustomAPIDestinationCredentials{
			CustomHeaders: NewTypedMap(map[string]TypedObject[CustomAPIDestinationCustomHeader]{
				key: NewTypedObject(CustomAPIDestinationCustomHeader{
					ValueWO:  types.StringValue("secret"),
					IsSecret: types.BoolValue(true),
				}),
			}),
		}),
	}

	_, values, diags := customAPIDestinationModelWithWriteOnlyCredentials(plan, config)
	require.False(t, diags.HasError(), diags)
	assert.Contains(t, values, expectedPath.String())
	assert.NotContains(t, values, `credentials.custom_headers["x-"secret"\key"].value_wo`)
}

func TestWriteOnlyCredentialVerifiersChangedDetectsRemovedWriteOnlyCredential(t *testing.T) {
	t.Parallel()

	private := fakePrivateStateReader{
		values: map[string][]byte{
			writeOnlyCredentialVerifiersPrivateKey: mustJSONMarshal(t, map[string]string{
				path.Root("credentials").AtName("client_key_wo").String(): mustWriteOnlyCredentialVerifier(t, path.Root("credentials").AtName("client_key_wo"), types.StringValue("previous-secret")),
			}),
		},
	}

	changed, diags := writeOnlyCredentialVerifiersChanged(context.Background(), private, writeOnlyCredentialValues{})

	require.False(t, diags.HasError(), diags)
	assert.True(t, changed)
}

func TestWriteOnlyCredentialVerifiersChangedDetectsAddedOrChangedWriteOnlyCredential(t *testing.T) {
	t.Parallel()

	apiKeyPath := path.Root("credentials").AtName("api_key_wo")
	clientKeyPath := path.Root("credentials").AtName("client_key_wo")

	private := fakePrivateStateReader{
		values: map[string][]byte{
			writeOnlyCredentialVerifiersPrivateKey: mustJSONMarshal(t, map[string]string{
				apiKeyPath.String(): mustWriteOnlyCredentialVerifier(t, apiKeyPath, types.StringValue("api-key")),
			}),
		},
	}

	changed, diags := writeOnlyCredentialVerifiersChanged(context.Background(), private, writeOnlyCredentialValues{
		apiKeyPath.String():    types.StringValue("api-key"),
		clientKeyPath.String(): types.StringValue("client-key"),
	})

	require.False(t, diags.HasError(), diags)
	assert.True(t, changed)
}

func TestWriteOnlyCredentialVerifiersChangedDetectsChangedWriteOnlyCredential(t *testing.T) {
	t.Parallel()

	apiKeyPath := path.Root("credentials").AtName("api_key_wo")

	private := fakePrivateStateReader{
		values: map[string][]byte{
			writeOnlyCredentialVerifiersPrivateKey: mustJSONMarshal(t, map[string]string{
				apiKeyPath.String(): mustWriteOnlyCredentialVerifier(t, apiKeyPath, types.StringValue("api-key")),
			}),
		},
	}

	changed, diags := writeOnlyCredentialVerifiersChanged(context.Background(), private, writeOnlyCredentialValues{
		apiKeyPath.String(): types.StringValue("rotated-api-key"),
	})

	require.False(t, diags.HasError(), diags)
	assert.True(t, changed)
}

func TestWriteOnlyCredentialVerifiersChangedIgnoresUnchangedValues(t *testing.T) {
	t.Parallel()

	apiKeyPath := path.Root("credentials").AtName("api_key_wo")
	values := writeOnlyCredentialValues{
		apiKeyPath.String(): types.StringValue("current-secret"),
	}

	private := fakePrivateStateReader{
		values: map[string][]byte{
			writeOnlyCredentialVerifiersPrivateKey: mustJSONMarshal(t, map[string]string{
				apiKeyPath.String(): mustWriteOnlyCredentialVerifier(t, apiKeyPath, values[apiKeyPath.String()]),
			}),
		},
	}

	changed, diags := writeOnlyCredentialVerifiersChanged(context.Background(), private, values)

	require.False(t, diags.HasError(), diags)
	assert.False(t, changed)
}

type fakePrivateStateReader struct {
	values map[string][]byte
}

func (r fakePrivateStateReader) GetKey(_ context.Context, key string) ([]byte, diag.Diagnostics) {
	return r.values[key], nil
}

func mustJSONMarshal(t *testing.T, value any) []byte {
	t.Helper()

	data, err := json.Marshal(value)
	require.NoError(t, err)

	return data
}

func mustWriteOnlyCredentialVerifier(t *testing.T, argumentPath path.Path, value types.String) string {
	t.Helper()

	verifier, err := writeOnlyCredentialVerifier(argumentPath.String(), value)
	require.NoError(t, err)

	return verifier
}
