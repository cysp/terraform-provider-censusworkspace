package provider_test

import (
	"testing"

	. "github.com/cysp/terraform-provider-censusworkspace/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestNewBigQueryDestinationConnectionDetailsFromResponse(t *testing.T) {
	t.Parallel()

	testcases := map[string]struct {
		input    []byte
		expected TypedObject[BigQueryDestinationConnectionDetails]
	}{
		"null": {
			input:    nil,
			expected: NewTypedObjectNull[BigQueryDestinationConnectionDetails](),
		},
		"empty": {
			input:    []byte(`{}`),
			expected: NewTypedObject(BigQueryDestinationConnectionDetails{}),
		},
		"project, location": {
			input: []byte(`{"project_id":"project-id","location":"location"}`),
			expected: NewTypedObject(BigQueryDestinationConnectionDetails{
				ProjectID: types.StringValue("project-id"),
				Location:  types.StringValue("location"),
			}),
		},
		"project, location, service_account_email": {
			input: []byte(`{"project_id":"project-id","location":"location","service_account_email":"service-account","service_account_key":null}`),
			expected: NewTypedObject(BigQueryDestinationConnectionDetails{
				ProjectID:           types.StringValue("project-id"),
				Location:            types.StringValue("location"),
				ServiceAccountEmail: types.StringValue("service-account"),
			}),
		},
		"project, location, service_account_email, service_account_key": {
			input: []byte(`{"project_id":"project-id","location":"location","service_account_email":"service-account","service_account_key":"service-account-key"}`),
			expected: NewTypedObject(BigQueryDestinationConnectionDetails{
				ProjectID:           types.StringValue("project-id"),
				Location:            types.StringValue("location"),
				ServiceAccountEmail: types.StringValue("service-account"),
				ServiceAccountKey:   types.StringValue("service-account-key"),
			}),
		},
	}

	for name, testcase := range testcases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			actual, diags := NewBigQueryDestinationConnectionDetailsFromResponse(t.Context(), path.Root("connection_details"), testcase.input)
			assert.Empty(t, diags)

			assert.Equal(t, testcase.expected, actual)
		})
	}
}
