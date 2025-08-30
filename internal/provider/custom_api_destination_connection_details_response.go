package provider

import (
	"github.com/go-faster/jx"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (cd *CustomAPIDestinationConnectionDetails) Decode(dec *jx.Decoder) error {
	//nolint:wrapcheck
	return dec.Obj(func(dec *jx.Decoder, key string) error {
		switch key {
		case "api_version":
			value, err := dec.Int64()
			if err != nil {
				//nolint:wrapcheck
				return err
			}

			cd.APIVersion = types.Int64Value(value)

		case "webhook_url":
			value, err := dec.Str()
			if err != nil {
				//nolint:wrapcheck
				return err
			}

			cd.WebhookURL = types.StringValue(value)

		case "custom_headers":
			if dec.Next() == jx.Null {
				return dec.Skip()
			}

			customHeaders := make(map[string]TypedObject[CustomAPIDestinationCustomHeader], 0)

			customHeadersDecodeError := dec.Obj(func(d *jx.Decoder, key string) error {
				customHeader := CustomAPIDestinationCustomHeader{}

				customHeaderDecodeErr := customHeader.Decode(d)
				if customHeaderDecodeErr != nil {
					return customHeaderDecodeErr
				}

				customHeaders[key] = NewTypedObject(customHeader)

				return nil
			})
			if customHeadersDecodeError != nil {
				//nolint:wrapcheck
				return customHeadersDecodeError
			}

			cd.CustomHeaders = NewTypedMap(customHeaders)

		default:
			return dec.Skip()
		}

		return nil
	})
}

func (ch *CustomAPIDestinationCustomHeader) Decode(dec *jx.Decoder) error {
	//nolint:wrapcheck
	return dec.Obj(func(dec *jx.Decoder, key string) error {
		switch key {
		case "value":
			if dec.Next() == jx.Null {
				err := dec.Null()
				if err != nil {
					//nolint:wrapcheck
					return err
				}

				ch.Value = types.StringNull()
			} else {
				value, err := dec.Str()
				if err != nil {
					//nolint:wrapcheck
					return err
				}

				ch.Value = types.StringValue(value)
			}

			return nil

		case "is_secret":
			value, err := dec.Bool()
			if err != nil {
				//nolint:wrapcheck
				return err
			}

			ch.IsSecret = types.BoolValue(value)

			return nil

		default:
			return dec.Skip()
		}
	})
}
