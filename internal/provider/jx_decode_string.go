package provider

import (
	"github.com/go-faster/jx"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func JxDecodeStringValue(dec *jx.Decoder, value *types.String) error {
	if dec.Next() == jx.Null {
		err := dec.Null()
		if err != nil {
			//nolint:wrapcheck
			return err
		}

		*value = types.StringNull()

		return nil
	}

	v, err := dec.Str()
	if err != nil {
		//nolint:wrapcheck
		return err
	}

	*value = types.StringValue(v)

	return nil
}
