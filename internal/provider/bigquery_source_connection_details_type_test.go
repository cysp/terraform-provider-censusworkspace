package provider_test

import (
	"testing"

	"github.com/cysp/terraform-provider-censusworkspace/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBigQuerySourceConnectionDetailsTypeValueFromObject(t *testing.T) {
	t.Parallel()

	ctx := t.Context()

	typ := provider.BigQuerySourceConnectionDetails{}.ObjectType(ctx)

	t.Run("unknown", func(t *testing.T) {
		t.Parallel()

		value := types.ObjectUnknown(typ.AttrTypes)

		object, diags := provider.BigQuerySourceConnectionDetailsType{}.ValueFromObject(ctx, value)

		assert.True(t, object.IsUnknown())
		assert.Empty(t, diags)
	})

	t.Run("null", func(t *testing.T) {
		t.Parallel()

		value := types.ObjectNull(typ.AttrTypes)

		object, diags := provider.BigQuerySourceConnectionDetailsType{}.ValueFromObject(ctx, value)

		assert.True(t, object.IsNull())
		assert.Empty(t, diags)
	})

	t.Run("value", func(t *testing.T) {
		t.Parallel()

		value, diags := types.ObjectValue(typ.AttrTypes, map[string]attr.Value{
			"project_id": types.StringValue("project-id"),
			"location":   types.StringValue("location"),
		})
		require.Empty(t, diags)
		require.False(t, diags.HasError())

		object, diags := provider.BigQuerySourceConnectionDetailsType{}.ValueFromObject(ctx, value)

		assert.False(t, diags.HasError())
		assert.False(t, object.IsNull())
		assert.False(t, object.IsUnknown())

		transformation, transformationOk := object.(provider.BigQuerySourceConnectionDetails)
		assert.True(t, transformationOk)
		assert.Equal(t, "project-id", transformation.ProjectID.ValueString())
		assert.Equal(t, "location", transformation.Location.ValueString())
		assert.True(t, transformation.ServiceAccount.IsNull())
	})
}
