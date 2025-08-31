package provider

import (
	"context"

	"github.com/go-faster/jx"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
)

func NewCustomAPIDestinationConnectionDetailsFromResponse(_ context.Context, path path.Path, data []byte) (TypedObject[CustomAPIDestinationConnectionDetails], diag.Diagnostics) {
	if data == nil {
		return NewTypedObjectNull[CustomAPIDestinationConnectionDetails](), nil
	}

	diags := diag.Diagnostics{}

	var connectionDetails CustomAPIDestinationConnectionDetails

	dec := jx.DecodeBytes(data)

	err := connectionDetails.Decode(dec)
	if err != nil {
		diags.AddAttributeError(path, "Error decoding value", err.Error())
	}

	return NewTypedObject(connectionDetails), diags
}

func (cd *CustomAPIDestinationConnectionDetails) Decode(dec *jx.Decoder) error {
	//nolint:wrapcheck
	return dec.Obj(func(dec *jx.Decoder, key string) error {
		switch key {
		case "api_version":
			return JxDecodeInt64Value(dec, &cd.APIVersion)

		case "webhook_url":
			return JxDecodeStringValue(dec, &cd.WebhookURL)

		case "custom_headers":
			return DecodeConnectionDetailsCustomHeadersMap(dec, &cd.CustomHeaders)

		default:
			return dec.Skip()
		}
	})
}

func DecodeConnectionDetailsCustomHeadersMap(dec *jx.Decoder, value *TypedMap[TypedObject[CustomAPIDestinationCustomHeader]]) error {
	if dec.Next() == jx.Null {
		err := dec.Null()
		if err != nil {
			//nolint:wrapcheck
			return err
		}

		*value = NewTypedMapNull[TypedObject[CustomAPIDestinationCustomHeader]]()

		return nil
	}

	customHeaders := make(map[string]TypedObject[CustomAPIDestinationCustomHeader], 0)

	customHeadersDecodeErr := dec.Obj(func(d *jx.Decoder, key string) error {
		customHeader := CustomAPIDestinationCustomHeader{}

		customHeaderDecodeErr := customHeader.Decode(d)
		if customHeaderDecodeErr != nil {
			return customHeaderDecodeErr
		}

		customHeaders[key] = NewTypedObject(customHeader)

		return nil
	})
	if customHeadersDecodeErr != nil {
		//nolint:wrapcheck
		return customHeadersDecodeErr
	}

	*value = NewTypedMap(customHeaders)

	return nil
}

func (ch *CustomAPIDestinationCustomHeader) Decode(dec *jx.Decoder) error {
	//nolint:wrapcheck
	return dec.Obj(func(dec *jx.Decoder, key string) error {
		switch key {
		case "value":
			return JxDecodeStringValue(dec, &ch.Value)

		case "is_secret":
			return JxDecodeBoolValue(dec, &ch.IsSecret)

		default:
			return dec.Skip()
		}
	})
}
