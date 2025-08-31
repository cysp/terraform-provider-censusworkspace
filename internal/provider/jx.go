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
