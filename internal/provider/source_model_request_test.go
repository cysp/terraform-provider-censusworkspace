package provider_test

import (
	"testing"

	. "github.com/cysp/terraform-provider-censusworkspace/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestSourceModelToCreateSourceDataSendsEmptyName(t *testing.T) {
	t.Parallel()

	model := SourceModel{}
	model.Name = types.StringValue("")
	model.Type = types.StringValue("big_query")
	model.Credentials = jsontypes.NewNormalizedNull()

	body, diags := model.ToCreateSourceData(t.Context())
	require.Empty(t, diags)

	name, ok := body.Connection.Name.Get()
	require.True(t, ok)
	require.Empty(t, name)
}

func TestSourceModelToUpdateSourceDataSendsEmptyName(t *testing.T) {
	t.Parallel()

	model := SourceModel{}
	model.Name = types.StringValue("")
	model.Credentials = jsontypes.NewNormalizedNull()

	body, diags := model.ToUpdateSourceData(t.Context())
	require.Empty(t, diags)

	name, ok := body.Connection.Name.Get()
	require.True(t, ok)
	require.Empty(t, name)
}

func TestBigQuerySourceModelToCreateSourceDataSendsEmptyName(t *testing.T) {
	t.Parallel()

	model := BigQuerySourceModel{}
	model.Name = types.StringValue("")
	model.Credentials = NewTypedObject(BigQuerySourceCredentials{
		ProjectID:         types.StringValue("project-id"),
		Location:          types.StringValue("US"),
		ServiceAccountKey: NewTypedObjectNull[BigQuerySourceCredentialsServiceAccountKey](),
	})

	body, diags := model.ToCreateSourceData(t.Context())
	require.Empty(t, diags)

	name, ok := body.Connection.Name.Get()
	require.True(t, ok)
	require.Empty(t, name)
}

func TestBigQuerySourceModelToUpdateSourceDataSendsEmptyName(t *testing.T) {
	t.Parallel()

	model := BigQuerySourceModel{}
	model.Name = types.StringValue("")
	model.Credentials = NewTypedObject(BigQuerySourceCredentials{
		ProjectID:         types.StringValue("project-id"),
		Location:          types.StringValue("US"),
		ServiceAccountKey: NewTypedObjectNull[BigQuerySourceCredentialsServiceAccountKey](),
	})

	body, diags := model.ToUpdateSourceData(t.Context())
	require.Empty(t, diags)

	name, ok := body.Connection.Name.Get()
	require.True(t, ok)
	require.Empty(t, name)
}
