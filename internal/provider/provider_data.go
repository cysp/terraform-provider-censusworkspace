package provider

import (
	cm "github.com/cysp/terraform-provider-censusworkspace/internal/census-management-go"
)

//nolint:revive
type ProviderData struct {
	client *cm.Client
}
