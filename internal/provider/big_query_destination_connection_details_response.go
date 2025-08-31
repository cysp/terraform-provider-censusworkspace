package provider

import (
	"context"

	"github.com/go-faster/jx"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
)

func NewBigQueryDestinationConnectionDetailsFromResponse(_ context.Context, path path.Path, data jx.Raw) (TypedObject[BigQueryDestinationConnectionDetails], diag.Diagnostics) {
	if data == nil {
		return NewTypedObjectNull[BigQueryDestinationConnectionDetails](), nil
	}

	diags := diag.Diagnostics{}

	var connectionDetails BigQueryDestinationConnectionDetails

	dec := jx.DecodeBytes(data)

	err := connectionDetails.Decode(dec)
	if err != nil {
		diags.AddAttributeError(path, "Error decoding value", err.Error())
	}

	return NewTypedObject(connectionDetails), diags
}

func (cd *BigQueryDestinationConnectionDetails) Decode(dec *jx.Decoder) error {
	//nolint:wrapcheck
	return dec.Obj(func(dec *jx.Decoder, key string) error {
		switch key {
		case "project_id":
			return JxDecodeStringValue(dec, &cd.ProjectID)

		case "location":
			return JxDecodeStringValue(dec, &cd.Location)

		case "service_account_email":
			return JxDecodeStringValue(dec, &cd.ServiceAccountEmail)

		case "service_account_key":
			return JxDecodeStringValue(dec, &cd.ServiceAccountKey)

		default:
			return dec.Skip()
		}
	})
}
