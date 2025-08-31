package provider_test

import (
	"testing"

	. "github.com/cysp/terraform-provider-censusworkspace/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestNewCustomAPIDestinationConnectionDetailsFromResponse(t *testing.T) {
	t.Parallel()

	testcases := map[string]struct {
		input    []byte
		expected TypedObject[CustomAPIDestinationConnectionDetails]
	}{
		"null": {
			input:    nil,
			expected: NewTypedObjectNull[CustomAPIDestinationConnectionDetails](),
		},
		"empty": {
			input:    []byte(`{}`),
			expected: NewTypedObject(CustomAPIDestinationConnectionDetails{}),
		},
		"webhook_url": {
			input: []byte(`{"api_version":1,"webhook_url":"https://example.org/census-destination"}`),
			expected: NewTypedObject(CustomAPIDestinationConnectionDetails{
				APIVersion: types.Int64Value(1),
				WebhookURL: types.StringValue("https://example.org/census-destination"),
			}),
		},
		"webhook_url, custom_headers=null": {
			input: []byte(`{"api_version":1,"webhook_url":"https://example.org/census-destination","custom_headers":null}`),
			expected: NewTypedObject(CustomAPIDestinationConnectionDetails{
				APIVersion:    types.Int64Value(1),
				WebhookURL:    types.StringValue("https://example.org/census-destination"),
				CustomHeaders: NewTypedMapNull[TypedObject[CustomAPIDestinationCustomHeader]](),
			}),
		},
		"webhook_url, custom_headers": {
			input: []byte(`{"api_version":1,"webhook_url":"https://example.org/census-destination","custom_headers":{"x-client-id":{"value":"client-id"}}}`),
			expected: NewTypedObject(CustomAPIDestinationConnectionDetails{
				APIVersion: types.Int64Value(1),
				WebhookURL: types.StringValue("https://example.org/census-destination"),
				CustomHeaders: NewTypedMap(map[string]TypedObject[CustomAPIDestinationCustomHeader]{
					"x-client-id": NewTypedObject(CustomAPIDestinationCustomHeader{
						Value:    types.StringValue("client-id"),
						IsSecret: types.BoolNull(),
					}),
				}),
			}),
		},
		"webhook_url, custom_headers=secret": {
			input: []byte(`{"api_version":1,"webhook_url":"https://example.org/census-destination","custom_headers":{"x-client-id":{"value":"client-id","is_secret":false},"x-client-secret":{"value":null,"is_secret":true}}}`),
			expected: NewTypedObject(CustomAPIDestinationConnectionDetails{
				APIVersion: types.Int64Value(1),
				WebhookURL: types.StringValue("https://example.org/census-destination"),
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
		},
	}

	for name, testcase := range testcases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			actual, diags := NewCustomAPIDestinationConnectionDetailsFromResponse(t.Context(), path.Root("connection_details"), testcase.input)
			assert.Empty(t, diags)

			assert.Equal(t, testcase.expected, actual)
		})
	}
}
