package provider

import (
	cm "github.com/cysp/terraform-provider-censusworkspace/internal/census-management-go"
)

type CensusProviderData struct {
	client *cm.Client
}
