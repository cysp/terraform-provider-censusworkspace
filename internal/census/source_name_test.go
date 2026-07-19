package census_test

import (
	"testing"

	"github.com/cysp/terraform-provider-censusworkspace/internal/census"
	"github.com/stretchr/testify/require"
)

func TestCanonicalizeSourceName(t *testing.T) {
	t.Parallel()

	tests := map[string]string{
		"bigquery - xxx-census-dev": "bigquery_xxx_census_dev",
		"source  -- -  name":        "source_name",
		"source_name":               "source_name",
		" - source - ":              "_source_",
	}

	for name, expected := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, expected, census.CanonicalizeSourceName(name))
		})
	}
}
