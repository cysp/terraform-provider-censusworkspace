package provider_test

import (
	"testing"

	. "github.com/cysp/terraform-provider-censusworkspace/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestNewBrazeDestinationConnectionDetailsFromResponse(t *testing.T) {
	t.Parallel()

	testcases := map[string]struct {
		input    []byte
		expected TypedObject[BrazeDestinationConnectionDetails]
	}{
		"null": {
			input:    nil,
			expected: NewTypedObjectNull[BrazeDestinationConnectionDetails](),
		},
		"empty": {
			input:    []byte(`{}`),
			expected: NewTypedObject(BrazeDestinationConnectionDetails{}),
		},
		"instance_url": {
			input: []byte(`{"instance_url":"instance-url"}`),
			expected: NewTypedObject(BrazeDestinationConnectionDetails{
				InstanceURL: types.StringValue("instance-url"),
			}),
		},
		"instance_url=null": {
			input: []byte(`{"instance_url":null}`),
			expected: NewTypedObject(BrazeDestinationConnectionDetails{
				InstanceURL: types.StringNull(),
			}),
		},
	}

	for name, testcase := range testcases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			actual, diags := NewBrazeDestinationConnectionDetailsFromResponse(t.Context(), path.Root("connection_details"), testcase.input)
			assert.Empty(t, diags)

			assert.Equal(t, testcase.expected, actual)
		})
	}
}
