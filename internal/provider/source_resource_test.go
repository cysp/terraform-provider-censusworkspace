package provider_test

import (
	"regexp"
	"testing"

	cm "github.com/cysp/terraform-provider-censusworkspace/internal/census-management-go"
	cmt "github.com/cysp/terraform-provider-censusworkspace/internal/census-management-go/testing"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

//nolint:paralleltest
func TestAccSourceResourceImport(t *testing.T) {
	server, err := cmt.NewCensusManagementServer()
	if err != nil {
		t.Fatal(err)
	}

	configVariables := config.Variables{
		"source_id": config.StringVariable("12345"),
	}

	server.Handler().Sources["12345"] = &cm.SourceData{
		ID:   12345,
		Name: "Test Source",
	}

	CensusProviderMockableResourceTest(t, server, resource.TestCase{
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: configVariables,
				ResourceName:    "census_source.test",
				ImportState:     true,
				ImportStateId:   "12345",
			},
		},
	})
}

//nolint:paralleltest
func TestAccSourceResourceImportNotFound(t *testing.T) {
	server, err := cmt.NewCensusManagementServer()
	if err != nil {
		t.Fatal(err)
	}

	configVariables := config.Variables{}

	CensusProviderMockableResourceTest(t, server, resource.TestCase{
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: configVariables,
				ResourceName:    "census_source.test",
				ImportState:     true,
				ImportStateId:   "99999",
				ExpectError:     regexp.MustCompile(`Cannot import non-existent remote object`),
			},
		},
	})
}

//nolint:paralleltest
func TestAccSourceResourceCreateUpdateDelete(t *testing.T) {
	server, err := cmt.NewCensusManagementServer()
	if err != nil {
		t.Fatal(err)
	}

	configVariables := config.Variables{}

	CensusProviderMockedResourceTest(t, server, resource.TestCase{
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: configVariables,
			},
			{
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: configVariables,
			},
		},
	})
}
