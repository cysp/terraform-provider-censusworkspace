//nolint:dupl
package provider_test

import (
	"testing"

	. "github.com/cysp/terraform-provider-censusworkspace/internal/provider"
	"github.com/go-faster/jx"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJxDecodeStringValue(t *testing.T) {
	t.Parallel()

	testcases := map[string]struct {
		input     string
		expected  types.String
		expectErr bool
	}{
		"empty": {
			input:     "",
			expectErr: true,
		},
		"null": {
			input:     "null",
			expected:  types.StringNull(),
			expectErr: false,
		},
		"valid": {
			input:     `"12345"`,
			expected:  types.StringValue("12345"),
			expectErr: false,
		},
		"invalid": {
			input:     "12345",
			expectErr: true,
		},
	}

	for name, testcase := range testcases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			dec := jx.DecodeStr(testcase.input)

			var actual types.String

			err := JxDecodeStringValue(dec, &actual)

			if testcase.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, testcase.expected, actual)
		})
	}
}
