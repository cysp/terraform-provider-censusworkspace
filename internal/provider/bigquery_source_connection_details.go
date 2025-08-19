package provider

import (
	"context"

	tpfr "github.com/cysp/terraform-provider-censusworkspace/internal/terraform-plugin-framework-reflection"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

type BigQuerySourceConnectionDetails struct {
	ProjectID      types.String `tfsdk:"project_id"`
	Location       types.String `tfsdk:"location"`
	ServiceAccount types.String `tfsdk:"service_account"`
	state          attr.ValueState
}

func NewBigQuerySourceConnectionDetailsNull() BigQuerySourceConnectionDetails {
	return BigQuerySourceConnectionDetails{
		state: attr.ValueStateNull,
	}
}

func NewBigQuerySourceConnectionDetailsUnknown() BigQuerySourceConnectionDetails {
	return BigQuerySourceConnectionDetails{
		state: attr.ValueStateUnknown,
	}
}

func NewBigQuerySourceConnectionDetailsKnown() BigQuerySourceConnectionDetails {
	return BigQuerySourceConnectionDetails{
		state: attr.ValueStateKnown,
	}
}

func NewBigQuerySourceConnectionDetailsKnownFromAttributes(ctx context.Context, attributes map[string]attr.Value) (BigQuerySourceConnectionDetails, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	value := BigQuerySourceConnectionDetails{
		state: attr.ValueStateKnown,
	}

	setAttributesDiags := tpfr.SetAttributesInValue(ctx, &value, attributes)
	diags = append(diags, setAttributesDiags...)

	return value, diags
}

//nolint:ireturn
func (v BigQuerySourceConnectionDetails) CustomType(ctx context.Context) basetypes.ObjectTypable {
	return BigQuerySourceConnectionDetailsType{ObjectType: v.ObjectType(ctx)}
}

var _ basetypes.ObjectValuable = BigQuerySourceConnectionDetails{}

//nolint:ireturn
func (v BigQuerySourceConnectionDetails) Type(ctx context.Context) attr.Type {
	return BigQuerySourceConnectionDetailsType{ObjectType: v.ObjectType(ctx)}
}

func (v BigQuerySourceConnectionDetails) ObjectType(ctx context.Context) types.ObjectType {
	return types.ObjectType{AttrTypes: v.ObjectAttrTypes(ctx)}
}

func (v BigQuerySourceConnectionDetails) ObjectAttrTypes(ctx context.Context) map[string]attr.Type {
	return ObjectAttrTypesFromSchemaAttributes(ctx, v.SchemaAttributes(ctx))
}

func (v BigQuerySourceConnectionDetails) Equal(o attr.Value) bool {
	return tpfr.ValuesEqual[BigQuerySourceConnectionDetails](v, o)
}

func (v BigQuerySourceConnectionDetails) IsNull() bool {
	return v.state == attr.ValueStateNull
}

func (v BigQuerySourceConnectionDetails) IsUnknown() bool {
	return v.state == attr.ValueStateUnknown
}

func (v BigQuerySourceConnectionDetails) String() string {
	return "BigQuerySourceConnectionDetails"
}

func (v BigQuerySourceConnectionDetails) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	//nolint:wrapcheck
	return tpfr.ValueToTerraformValue(ctx, v, v.state)
}

func (v BigQuerySourceConnectionDetails) ToObjectValue(ctx context.Context) (types.Object, diag.Diagnostics) {
	return tpfr.ValueToObjectValue(ctx, v)
}
