package provider

import (
	"context"

	"github.com/go-faster/jx"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewBigQuerySourceConnectionDetailsFromResponse(ctx context.Context, path path.Path, data jx.Raw) (BigQuerySourceConnectionDetails, diag.Diagnostics) {
	if data == nil {
		return NewBigQuerySourceConnectionDetailsNull(), nil
	}

	diags := diag.Diagnostics{}

	values := make(map[string]attr.Value)

	dec := jx.DecodeBytes(data)
	decodeErr := dec.Obj(func(dec *jx.Decoder, key string) error {
		switch key {
		case "project_id", "location", "service_account":
			value, err := dec.Str()
			if err != nil {
				return err
			}

			values[key] = types.StringValue(value)
		}

		return nil
	})

	if decodeErr != nil {
		diags.AddAttributeError(path, "Error decoding connection details", decodeErr.Error())
	}

	value, valueDiags := NewBigQuerySourceConnectionDetailsKnownFromAttributes(ctx, values)
	diags.Append(valueDiags...)

	return value, diags
}
