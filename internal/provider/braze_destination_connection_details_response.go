package provider

import (
	"context"

	"github.com/go-faster/jx"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
)

func NewBrazeDestinationConnectionDetailsFromResponse(_ context.Context, path path.Path, data jx.Raw) (TypedObject[BrazeDestinationConnectionDetails], diag.Diagnostics) {
	if data == nil {
		return NewTypedObjectNull[BrazeDestinationConnectionDetails](), nil
	}

	diags := diag.Diagnostics{}

	var connectionDetails BrazeDestinationConnectionDetails

	dec := jx.DecodeBytes(data)

	err := connectionDetails.Decode(dec)
	if err != nil {
		diags.AddAttributeError(path, "Error decoding value", err.Error())
	}

	return NewTypedObject(connectionDetails), diags
}

func (cd *BrazeDestinationConnectionDetails) Decode(dec *jx.Decoder) error {
	//nolint:wrapcheck
	return dec.Obj(func(dec *jx.Decoder, key string) error {
		switch key {
		case "instance_url":
			return JxDecodeStringValue(dec, &cd.InstanceURL)

		default:
			return dec.Skip()
		}
	})
}
