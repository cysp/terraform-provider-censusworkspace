package provider

import (
	"context"

	"github.com/go-faster/jx"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
)

func NewBigQuerySourceConnectionDetailsFromResponse(_ context.Context, path path.Path, data jx.Raw) (TypedObject[BigQuerySourceConnectionDetails], diag.Diagnostics) {
	if data == nil {
		return NewTypedObjectNull[BigQuerySourceConnectionDetails](), nil
	}

	diags := diag.Diagnostics{}

	var connectionDetails BigQuerySourceConnectionDetails

	dec := jx.DecodeBytes(data)

	err := connectionDetails.Decode(dec)
	if err != nil {
		diags.AddAttributeError(path, "Error decoding value", err.Error())
	}

	return NewTypedObject(connectionDetails), diags
}

func (cd *BigQuerySourceConnectionDetails) Decode(dec *jx.Decoder) error {
	//nolint:wrapcheck
	return dec.Obj(func(dec *jx.Decoder, key string) error {
		switch key {
		case "project_id":
			return JxDecodeStringValue(dec, &cd.ProjectID)

		case "location":
			return JxDecodeStringValue(dec, &cd.Location)

		case "service_account":
			return JxDecodeStringValue(dec, &cd.ServiceAccount)

		default:
			return dec.Skip()
		}
	})
}
