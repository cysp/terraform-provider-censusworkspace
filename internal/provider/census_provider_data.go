package provider

import (
	cm "github.com/cysp/terraform-provider-censusworkspace/internal/census-management-go/client"
)

type CensusProviderData struct {
	client *cm.Client
}
