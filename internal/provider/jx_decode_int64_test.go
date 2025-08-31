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

func TestJxDecodeInt64Value(t *testing.T) {
	t.Parallel()

	testcases := map[string]struct {
		input     string
		expected  types.Int64
		expectErr bool
	}{
		"empty": {
			input:     "",
			expectErr: true,
		},
		"null": {
			input:     "null",
			expected:  types.Int64Null(),
			expectErr: false,
		},
		"valid": {
			input:     "12345",
			expected:  types.Int64Value(12345),
			expectErr: false,
		},
		"invalid": {
			input:     `"string"`,
			expectErr: true,
		},
	}

	for name, testcase := range testcases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			dec := jx.DecodeStr(testcase.input)

			var actual types.Int64

			err := JxDecodeInt64Value(dec, &actual)

			if testcase.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, testcase.expected, actual)
		})
	}
}
