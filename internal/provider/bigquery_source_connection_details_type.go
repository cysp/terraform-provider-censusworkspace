//nolint:dupl
package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

type BigQuerySourceConnectionDetailsType struct {
	basetypes.ObjectType
}

var _ basetypes.ObjectTypable = BigQuerySourceConnectionDetailsType{}

func (t BigQuerySourceConnectionDetailsType) Equal(o attr.Type) bool {
	other, ok := o.(BigQuerySourceConnectionDetailsType)
	if !ok {
		return false
	}

	return t.ObjectType.Equal(other.ObjectType)
}

func (t BigQuerySourceConnectionDetailsType) String() string {
	return "BigQuerySourceConnectionDetailsType"
}

//nolint:ireturn
func (t BigQuerySourceConnectionDetailsType) TerraformType(ctx context.Context) tftypes.Type {
	return tftypes.Object{
		AttributeTypes: ObjectAttrTypesToTerraformTypes(ctx, BigQuerySourceConnectionDetails{}.ObjectAttrTypes(ctx)),
	}
}

//nolint:ireturn
func (t BigQuerySourceConnectionDetailsType) ValueFromTerraform(ctx context.Context, value tftypes.Value) (attr.Value, error) {
	if value.Type() == nil {
		return NewBigQuerySourceConnectionDetailsNull(), nil
	}

	if !value.Type().Equal(t.TerraformType(ctx)) {
		return nil, UnexpectedTerraformTypeError{Expected: t.TerraformType(ctx), Actual: value.Type()}
	}

	if value.IsNull() {
		return NewBigQuerySourceConnectionDetailsNull(), nil
	}

	if !value.IsKnown() {
		return NewBigQuerySourceConnectionDetailsUnknown(), nil
	}

	attributes, err := AttributesFromTerraformValue(ctx, t.AttrTypes, value)
	if err != nil {
		return nil, fmt.Errorf("failed to create BigQuerySourceConnectionDetails from Terraform: %w", err)
	}

	v, diags := NewBigQuerySourceConnectionDetailsKnownFromAttributes(ctx, attributes)

	return v, ErrorFromDiags(diags)
}

//nolint:ireturn
func (t BigQuerySourceConnectionDetailsType) ValueFromObject(ctx context.Context, value basetypes.ObjectValue) (basetypes.ObjectValuable, diag.Diagnostics) {
	switch {
	case value.IsNull():
		return NewBigQuerySourceConnectionDetailsNull(), nil
	case value.IsUnknown():
		return NewBigQuerySourceConnectionDetailsUnknown(), nil
	}

	return NewBigQuerySourceConnectionDetailsKnownFromAttributes(ctx, value.Attributes())
}
