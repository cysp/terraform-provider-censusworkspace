package provider_test

import (
	"testing"

	. "github.com/cysp/terraform-provider-censusworkspace/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestNewBigQuerySourceConnectionDetailsFromResponse(t *testing.T) {
	t.Parallel()

	testcases := map[string]struct {
		input    []byte
		expected TypedObject[BigQuerySourceConnectionDetails]
	}{
		"null": {
			input:    nil,
			expected: NewTypedObjectNull[BigQuerySourceConnectionDetails](),
		},
		"empty": {
			input:    []byte(`{}`),
			expected: NewTypedObject(BigQuerySourceConnectionDetails{}),
		},
		"project, location": {
			input: []byte(`{"project_id":"project-id","location":"location"}`),
			expected: NewTypedObject(BigQuerySourceConnectionDetails{
				ProjectID: types.StringValue("project-id"),
				Location:  types.StringValue("location"),
			}),
		},
		"project, location, service_account": {
			input: []byte(`{"project_id":"project-id","location":"location","service_account":"service-account"}`),
			expected: NewTypedObject(BigQuerySourceConnectionDetails{
				ProjectID:      types.StringValue("project-id"),
				Location:       types.StringValue("location"),
				ServiceAccount: types.StringValue("service-account"),
			}),
		},
		"project, location, service_account=null": {
			input: []byte(`{"project_id":"project-id","location":"location","service_account":null}`),
			expected: NewTypedObject(BigQuerySourceConnectionDetails{
				ProjectID: types.StringValue("project-id"),
				Location:  types.StringValue("location"),
			}),
		},
	}

	for name, testcase := range testcases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			actual, diags := NewBigQuerySourceConnectionDetailsFromResponse(t.Context(), path.Root("connection_details"), testcase.input)
			assert.Empty(t, diags)

			assert.Equal(t, testcase.expected, actual)
		})
	}
}
