package provider

import (
	"github.com/go-faster/jx"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func JxDecodeBoolValue(dec *jx.Decoder, value *types.Bool) error {
	if dec.Next() == jx.Null {
		err := dec.Null()
		if err != nil {
			//nolint:wrapcheck
			return err
		}

		*value = types.BoolNull()

		return nil
	}

	v, err := dec.Bool()
	if err != nil {
		//nolint:wrapcheck
		return err
	}

	*value = types.BoolValue(v)

	return nil
}
