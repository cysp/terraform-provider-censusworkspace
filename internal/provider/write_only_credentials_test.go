package provider //nolint:testpackage

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"strings"
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
	assert.Contains(t, values, mustWriteOnlyCredentialPath(t, path.Root("credentials").AtName("api_key_wo")))
	assert.Equal(t, "api-key", requestModel.Credentials.Value().APIKey.ValueString())
	assert.Equal(t, "client-key", requestModel.Credentials.Value().ClientKey.ValueString())

	stateModel := sanitizedBrazeDestinationCredentials(plan, plan, config, BrazeDestinationConnectionDetails{
		InstanceURL: types.StringValue("normalized-instance-url"),
	})
	assert.True(t, stateModel.Credentials.Value().APIKey.IsNull())
	assert.Equal(t, "client-key", stateModel.Credentials.Value().ClientKey.ValueString())
	assert.Equal(t, "normalized-instance-url", stateModel.Credentials.Value().InstanceURL.ValueString())
}

func TestBrazeDestinationWriteOnlyCredentialResolverUsesPlanWhenConfigCredentialsUnknown(t *testing.T) {
	t.Parallel()

	plan := BrazeDestinationModel{
		Credentials: NewTypedObject(BrazeDestinationCredentials{
			InstanceURL: types.StringValue("instance-url"),
			APIKey:      types.StringValue("api-key"),
			ClientKey:   types.StringValue("client-key"),
		}),
	}
	config := BrazeDestinationModel{
		Credentials: NewTypedObjectUnknown[BrazeDestinationCredentials](),
	}

	requestModel, values, diags := brazeDestinationModelWithWriteOnlyCredentials(plan, config)
	require.False(t, diags.HasError(), diags)
	assert.Empty(t, values)
	assert.Equal(t, "instance-url", requestModel.Credentials.Value().InstanceURL.ValueString())
	assert.Equal(t, "api-key", requestModel.Credentials.Value().APIKey.ValueString())
	assert.Equal(t, "client-key", requestModel.Credentials.Value().ClientKey.ValueString())
}

func TestBrazeDestinationWriteOnlyCredentialResolverTracksUnknownWriteOnlyValue(t *testing.T) {
	t.Parallel()

	writeOnlyPath := path.Root("credentials").AtName("api_key_wo")
	plan := BrazeDestinationModel{
		Credentials: NewTypedObject(BrazeDestinationCredentials{
			InstanceURL: types.StringValue("instance-url"),
			APIKey:      types.StringNull(),
		}),
	}
	config := BrazeDestinationModel{
		Credentials: NewTypedObject(BrazeDestinationCredentials{
			InstanceURL: types.StringValue("instance-url"),
			APIKeyWO:    types.StringUnknown(),
		}),
	}

	_, values, diags := brazeDestinationModelWithWriteOnlyCredentials(plan, config)

	require.False(t, diags.HasError(), diags)
	writeOnlyKey := mustWriteOnlyCredentialPath(t, writeOnlyPath)
	require.Contains(t, values, writeOnlyKey)
	assert.True(t, values[writeOnlyKey].IsUnknown())
}

func TestBrazeDestinationWriteOnlyCredentialResolverRejectsUnknownWriteOnlyValueForApply(t *testing.T) {
	t.Parallel()

	plan := BrazeDestinationModel{
		Credentials: NewTypedObject(BrazeDestinationCredentials{
			InstanceURL: types.StringValue("instance-url"),
			APIKey:      types.StringNull(),
		}),
	}
	config := BrazeDestinationModel{
		Credentials: NewTypedObject(BrazeDestinationCredentials{
			InstanceURL: types.StringValue("instance-url"),
			APIKeyWO:    types.StringUnknown(),
		}),
	}

	requestModel, values, diags := brazeDestinationModelWithWriteOnlyCredentials(plan, config)
	require.False(t, diags.HasError(), diags)
	assert.True(t, requestModel.Credentials.Value().APIKey.IsUnknown())

	diags = validateWriteOnlyCredentialValuesKnown(values)
	assert.True(t, diags.HasError())
}

func TestSanitizedBrazeDestinationCredentialsTreatsUnknownWriteOnlyValueAsConfigured(t *testing.T) {
	t.Parallel()

	plan := BrazeDestinationModel{
		Credentials: NewTypedObject(BrazeDestinationCredentials{
			InstanceURL: types.StringValue("instance-url"),
			APIKey:      types.StringUnknown(),
			ClientKey:   types.StringValue("client-key"),
		}),
	}
	config := BrazeDestinationModel{
		Credentials: NewTypedObject(BrazeDestinationCredentials{
			InstanceURL: types.StringValue("instance-url"),
			APIKeyWO:    types.StringUnknown(),
		}),
	}

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
	assert.Contains(t, values, mustWriteOnlyCredentialPath(t, path.Root("credentials").AtName("custom_headers").AtMapKey("x-client-secret").AtName("value_wo")))

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
	apiKeyKey := mustWriteOnlyCredentialPath(t, apiKeyPath)
	clientKeyKey := mustWriteOnlyCredentialPath(t, clientKeyPath)
	apiKeyVerifier, err := writeOnlyCredentialVerifier(apiKeyKey, value)
	require.NoError(t, err)

	assert.True(t, writeOnlyCredentialVerifierMatches(apiKeyKey, value, apiKeyVerifier))
	assert.False(t, writeOnlyCredentialVerifierMatches(clientKeyKey, value, apiKeyVerifier))
}

func TestResolveRequiredStringCredentialAllowsUnknownValues(t *testing.T) {
	t.Parallel()

	_, _, diags := resolveStringCredential(
		types.StringUnknown(),
		types.StringNull(),
		path.Root("credentials").AtName("api_key"),
		path.Root("credentials").AtName("api_key_wo"),
		true,
	)

	assert.False(t, diags.HasError(), diags)
}

func TestResolveRequiredStringCredentialRejectsMissingValues(t *testing.T) {
	t.Parallel()

	_, _, diags := resolveStringCredential(
		types.StringNull(),
		types.StringNull(),
		path.Root("credentials").AtName("api_key"),
		path.Root("credentials").AtName("api_key_wo"),
		true,
	)

	assert.True(t, diags.HasError())
}

func TestWriteOnlyCredentialVerifierRejectsUnexpectedParameters(t *testing.T) {
	t.Parallel()

	value := types.StringValue("shared-secret")
	apiKeyPath := path.Root("credentials").AtName("api_key_wo")
	apiKeyKey := mustWriteOnlyCredentialPath(t, apiKeyPath)
	verifier := mustWriteOnlyCredentialVerifier(t, apiKeyPath, value)

	assert.False(t, writeOnlyCredentialVerifierMatches(apiKeyKey, value, strings.Replace(verifier, "m=16384,t=1,p=1", "m=32768,t=1,p=1", 1)))
	assert.False(t, writeOnlyCredentialVerifierMatches(apiKeyKey, value, strings.Replace(verifier, "m=16384,t=1,p=1", "m=16384,t=2,p=1", 1)))
	assert.False(t, writeOnlyCredentialVerifierMatches(apiKeyKey, value, strings.Replace(verifier, "m=16384,t=1,p=1", "m=16384,t=1,p=2", 1)))

	parts := strings.Split(verifier, "$")
	require.Len(t, parts, 6)
	parts[4] = base64.RawStdEncoding.EncodeToString(make([]byte, writeOnlyCredentialSaltLength+1))
	assert.False(t, writeOnlyCredentialVerifierMatches(apiKeyKey, value, strings.Join(parts, "$")))
}

func TestWriteOnlyCredentialPathEncodesAttributePath(t *testing.T) {
	t.Parallel()

	key := mustWriteOnlyCredentialPath(t, path.Root("credentials").AtName("api_key_wo"))

	assert.JSONEq(t, `[{"type":"attribute","value":"credentials"},{"type":"attribute","value":"api_key_wo"}]`, string(key))
}

func TestWriteOnlyCredentialPathEncodesMapKeyPath(t *testing.T) {
	t.Parallel()

	key := mustWriteOnlyCredentialPath(t, path.Root("credentials").AtName("custom_headers").AtMapKey("x-client-secret").AtName("value_wo"))

	assert.JSONEq(t, `[{"type":"attribute","value":"credentials"},{"type":"attribute","value":"custom_headers"},{"type":"map_key","value":"x-client-secret"},{"type":"attribute","value":"value_wo"}]`, string(key))
}

func TestWriteOnlyCredentialPathMapKeysDoNotCollideWithDisplayPathSyntax(t *testing.T) {
	t.Parallel()

	key := `x-"secret"\key`
	expectedPath := path.Root("credentials").AtName("custom_headers").AtMapKey(key).AtName("value_wo")
	keyWithSpecialCharacters := mustWriteOnlyCredentialPath(t, expectedPath)
	displayLikeKey := mustWriteOnlyCredentialPath(t, path.Root("credentials").AtName("custom_headers").AtMapKey(`x-\"secret\"\\key`).AtName("value_wo"))

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
	assert.Contains(t, values, keyWithSpecialCharacters)
	assert.NotContains(t, values, `credentials.custom_headers["x-"secret"\key"].value_wo`)
	assert.NotEqual(t, keyWithSpecialCharacters, displayLikeKey)
}

func TestWriteOnlyCredentialPathRejectsUnsupportedPathSteps(t *testing.T) {
	t.Parallel()

	_, diags := writeOnlyCredentialPathFromPath(path.Root("credentials").AtListIndex(0))

	assert.True(t, diags.HasError())
}

func TestWriteOnlyCredentialValuesTracksOnlyConfiguredValues(t *testing.T) {
	t.Parallel()

	nullPath := path.Root("credentials").AtName("null_wo")
	unknownPath := path.Root("credentials").AtName("unknown_wo")
	knownPath := path.Root("credentials").AtName("known_wo")
	values := writeOnlyCredentialValues{}

	require.False(t, values.Add(nullPath, types.StringNull()).HasError())
	require.False(t, values.Add(unknownPath, types.StringUnknown()).HasError())
	require.False(t, values.Add(knownPath, types.StringValue("secret")).HasError())

	assert.NotContains(t, values, mustWriteOnlyCredentialPath(t, nullPath))
	assert.True(t, values[mustWriteOnlyCredentialPath(t, unknownPath)].IsUnknown())
	assert.Equal(t, "secret", values[mustWriteOnlyCredentialPath(t, knownPath)].ValueString())
}

func TestWriteOnlyCredentialVerifiersChangedDetectsRemovedWriteOnlyCredential(t *testing.T) {
	t.Parallel()

	private := fakePrivateStateReader{
		values: map[string][]byte{
			writeOnlyCredentialVerifiersPrivateKey: mustJSONMarshal(t, map[string]string{
				string(mustWriteOnlyCredentialPath(t, path.Root("credentials").AtName("client_key_wo"))): mustWriteOnlyCredentialVerifier(t, path.Root("credentials").AtName("client_key_wo"), types.StringValue("previous-secret")),
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
				string(mustWriteOnlyCredentialPath(t, apiKeyPath)): mustWriteOnlyCredentialVerifier(t, apiKeyPath, types.StringValue("api-key")),
			}),
		},
	}

	changed, diags := writeOnlyCredentialVerifiersChanged(context.Background(), private, writeOnlyCredentialValues{
		mustWriteOnlyCredentialPath(t, apiKeyPath):    types.StringValue("api-key"),
		mustWriteOnlyCredentialPath(t, clientKeyPath): types.StringValue("client-key"),
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
				string(mustWriteOnlyCredentialPath(t, apiKeyPath)): mustWriteOnlyCredentialVerifier(t, apiKeyPath, types.StringValue("api-key")),
			}),
		},
	}

	changed, diags := writeOnlyCredentialVerifiersChanged(context.Background(), private, writeOnlyCredentialValues{
		mustWriteOnlyCredentialPath(t, apiKeyPath): types.StringValue("rotated-api-key"),
	})

	require.False(t, diags.HasError(), diags)
	assert.True(t, changed)
}

func TestWriteOnlyCredentialVerifiersChangedIgnoresUnchangedValues(t *testing.T) {
	t.Parallel()

	apiKeyPath := path.Root("credentials").AtName("api_key_wo")
	apiKeyKey := mustWriteOnlyCredentialPath(t, apiKeyPath)
	values := writeOnlyCredentialValues{
		apiKeyKey: types.StringValue("current-secret"),
	}

	private := fakePrivateStateReader{
		values: map[string][]byte{
			writeOnlyCredentialVerifiersPrivateKey: mustJSONMarshal(t, map[string]string{
				string(apiKeyKey): mustWriteOnlyCredentialVerifier(t, apiKeyPath, values[apiKeyKey]),
			}),
		},
	}

	changed, diags := writeOnlyCredentialVerifiersChanged(context.Background(), private, values)

	require.False(t, diags.HasError(), diags)
	assert.False(t, changed)
}

func TestWriteOnlyCredentialVerifiersChangedDetectsUnknownConfiguredValue(t *testing.T) {
	t.Parallel()

	apiKeyPath := path.Root("credentials").AtName("api_key_wo")
	apiKeyKey := mustWriteOnlyCredentialPath(t, apiKeyPath)

	private := fakePrivateStateReader{
		values: map[string][]byte{
			writeOnlyCredentialVerifiersPrivateKey: mustJSONMarshal(t, map[string]string{
				string(apiKeyKey): mustWriteOnlyCredentialVerifier(t, apiKeyPath, types.StringValue("")),
			}),
		},
	}

	changed, diags := writeOnlyCredentialVerifiersChanged(context.Background(), private, writeOnlyCredentialValues{
		apiKeyKey: types.StringUnknown(),
	})

	require.False(t, diags.HasError(), diags)
	assert.True(t, changed)
}

func TestWriteOnlyCredentialVerifiersRejectUnknownConfiguredValue(t *testing.T) {
	t.Parallel()

	writer := fakePrivateStateWriter{}

	diags := writeWriteOnlyCredentialVerifiers(context.Background(), &writer, writeOnlyCredentialValues{
		mustWriteOnlyCredentialPath(t, path.Root("credentials").AtName("api_key_wo")): types.StringUnknown(),
	})

	assert.True(t, diags.HasError())
	assert.Empty(t, writer.values)
}

func TestWriteOnlyCredentialVerifiersWriteKnownConfiguredValue(t *testing.T) {
	t.Parallel()

	writer := fakePrivateStateWriter{}

	diags := writeWriteOnlyCredentialVerifiers(context.Background(), &writer, writeOnlyCredentialValues{
		mustWriteOnlyCredentialPath(t, path.Root("credentials").AtName("api_key_wo")): types.StringValue("secret"),
	})

	require.False(t, diags.HasError(), diags)
	assert.NotEmpty(t, writer.values[writeOnlyCredentialVerifiersPrivateKey])
}

func TestValidateWriteOnlyCredentialValuesKnownRejectsUnknownValue(t *testing.T) {
	t.Parallel()

	diags := validateWriteOnlyCredentialValuesKnown(writeOnlyCredentialValues{
		mustWriteOnlyCredentialPath(t, path.Root("credentials").AtName("api_key_wo")): types.StringUnknown(),
	})

	assert.True(t, diags.HasError())
}

type fakePrivateStateReader struct {
	values map[string][]byte
}

func (r fakePrivateStateReader) GetKey(_ context.Context, key string) ([]byte, diag.Diagnostics) {
	return r.values[key], nil
}

type fakePrivateStateWriter struct {
	values map[string][]byte
}

func (w *fakePrivateStateWriter) SetKey(_ context.Context, key string, value []byte) diag.Diagnostics {
	if w.values == nil {
		w.values = map[string][]byte{}
	}

	w.values[key] = value

	return nil
}

func mustJSONMarshal(t *testing.T, value any) []byte {
	t.Helper()

	data, err := json.Marshal(value)
	require.NoError(t, err)

	return data
}

func mustWriteOnlyCredentialVerifier(t *testing.T, argumentPath path.Path, value types.String) string {
	t.Helper()

	key := mustWriteOnlyCredentialPath(t, argumentPath)
	verifier, err := writeOnlyCredentialVerifier(key, value)
	require.NoError(t, err)

	return verifier
}

func mustWriteOnlyCredentialPath(t *testing.T, argumentPath path.Path) writeOnlyCredentialPath {
	t.Helper()

	key, diags := writeOnlyCredentialPathFromPath(argumentPath)
	require.False(t, diags.HasError(), diags)

	return key
}
