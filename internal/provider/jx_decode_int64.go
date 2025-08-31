package provider

import (
	"github.com/go-faster/jx"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func JxDecodeInt64Value(dec *jx.Decoder, value *types.Int64) error {
	if dec.Next() == jx.Null {
		err := dec.Null()
		if err != nil {
			//nolint:wrapcheck
			return err
		}

		*value = types.Int64Null()

		return nil
	}

	v, err := dec.Int64()
	if err != nil {
		//nolint:wrapcheck
		return err
	}

	*value = types.Int64Value(v)

	return nil
}
