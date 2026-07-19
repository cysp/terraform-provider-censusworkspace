package testing_test

import (
	"testing"

	cm "github.com/cysp/terraform-provider-censusworkspace/internal/census-management-go"
	cmt "github.com/cysp/terraform-provider-censusworkspace/internal/census-management-go/testing"
	"github.com/stretchr/testify/require"
)

func TestUpdateSourceWithUpdateSourceBodyCanonicalizesName(t *testing.T) {
	t.Parallel()

	source := cm.SourceData{
		Name: "bigquery_xxx_census_dev",
	}
	body := cm.UpdateSourceBody{}
	body.Connection.Name.SetTo("bigquery - xxx-census-dev")

	cmt.UpdateSourceWithUpdateSourceBody(&source, body)

	require.Equal(t, "bigquery_xxx_census_dev", source.Name)
}
