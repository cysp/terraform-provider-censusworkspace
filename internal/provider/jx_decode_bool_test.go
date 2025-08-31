package provider_test

import (
	"testing"

	. "github.com/cysp/terraform-provider-censusworkspace/internal/provider"
	"github.com/go-faster/jx"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJxDecodeBoolValue(t *testing.T) {
	t.Parallel()

	testcases := map[string]struct {
		input     string
		expected  types.Bool
		expectErr bool
	}{
		"empty": {
			input:     "",
			expectErr: true,
		},
		"null": {
			input:    "null",
			expected: types.BoolNull(),
		},
		"true": {
			input:    "true",
			expected: types.BoolValue(true),
		},
		"false": {
			input:    "false",
			expected: types.BoolValue(false),
		},
		"non-bool": {
			input:     `"string"`,
			expectErr: true,
		},
	}

	for name, testcase := range testcases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			dec := jx.DecodeStr(testcase.input)

			var actual types.Bool

			err := JxDecodeBoolValue(dec, &actual)

			if testcase.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, testcase.expected, actual)
		})
	}
}
