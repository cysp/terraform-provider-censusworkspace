package provider_test

import (
	"regexp"
	"strconv"
	"testing"

	cm "github.com/cysp/terraform-provider-censusworkspace/internal/census-management-go"
	cmt "github.com/cysp/terraform-provider-censusworkspace/internal/census-management-go/testing"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/stretchr/testify/require"
)

//nolint:paralleltest
func TestAccSQLDatasetResourceImport(t *testing.T) {
	server, err := cmt.NewCensusManagementServer()
	require.NoError(t, err)

	testDatasetID := int64(12345)
	testDatasetIDString := strconv.FormatInt(testDatasetID, 10)

	configVariables := config.Variables{
		"dataset_id":        config.StringVariable(testDatasetIDString),
		"dataset_name":      config.StringVariable("Test SQL Dataset"),
		"dataset_source_id": config.StringVariable("1"),
		"dataset_query":     config.StringVariable("SELECT 1"),
	}

	dataset := cm.DatasetData{}
	dataset.SetSQLDatasetData(cm.SQLDatasetData{
		ID:       testDatasetID,
		Name:     "Test SQL Dataset",
		Type:     cm.SQLDatasetDataTypeSQL,
		SourceID: 1,
		Query:    "SELECT 1",
	})

	server.Handler().Datasets[testDatasetIDString] = &dataset

	ProviderMockedResourceTest(t, server, resource.TestCase{
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: configVariables,
				ResourceName:    "censusworkspace_sql_dataset.test",
				ImportState:     true,
				ImportStateId:   testDatasetIDString,
			},
		},
	})
}

//nolint:paralleltest
func TestAccSQLDatasetResourceImportNotFound(t *testing.T) {
	server, err := cmt.NewCensusManagementServer()
	require.NoError(t, err)

	ProviderMockedResourceTest(t, server, resource.TestCase{
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestNameDirectory(),
				ResourceName:    "censusworkspace_sql_dataset.test",
				ImportState:     true,
				ImportStateId:   "99999",
				ExpectError:     regexp.MustCompile(`Cannot import non-existent remote object`),
			},
		},
	})
}

//nolint:paralleltest
func TestAccSQLDatasetResourceCreateUpdateDelete(t *testing.T) {
	server, err := cmt.NewCensusManagementServer()
	require.NoError(t, err)

	ProviderMockedResourceTest(t, server, resource.TestCase{
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"dataset_name":      config.StringVariable("Test SQL Dataset"),
					"dataset_source_id": config.StringVariable("1"),
					"dataset_query":     config.StringVariable("SELECT 1"),
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("censusworkspace_sql_dataset.test", plancheck.ResourceActionCreate),
						plancheck.ExpectUnknownValue("censusworkspace_sql_dataset.test", tfjsonpath.New("id")),
					},
					PostApplyPreRefresh: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue("censusworkspace_sql_dataset.test", tfjsonpath.New("id"), knownvalue.NotNull()),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue("censusworkspace_sql_dataset.test", tfjsonpath.New("id"), knownvalue.NotNull()),
						plancheck.ExpectKnownValue("censusworkspace_sql_dataset.test", tfjsonpath.New("name"), knownvalue.StringExact("Test SQL Dataset")),
						plancheck.ExpectKnownValue("censusworkspace_sql_dataset.test", tfjsonpath.New("query"), knownvalue.StringExact("SELECT 1")),
					},
				},
			},
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"dataset_name":      config.StringVariable("Test SQL Dataset (updated)"),
					"dataset_source_id": config.StringVariable("1"),
					"dataset_query":     config.StringVariable("SELECT 1"),
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("censusworkspace_sql_dataset.test", plancheck.ResourceActionUpdate),
						plancheck.ExpectKnownValue("censusworkspace_sql_dataset.test", tfjsonpath.New("id"), knownvalue.NotNull()),
						plancheck.ExpectKnownValue("censusworkspace_sql_dataset.test", tfjsonpath.New("name"), knownvalue.StringExact("Test SQL Dataset (updated)")),
						plancheck.ExpectKnownValue("censusworkspace_sql_dataset.test", tfjsonpath.New("query"), knownvalue.StringExact("SELECT 1")),
					},
				},
			},
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"dataset_name":      config.StringVariable("Test SQL Dataset (replaced)"),
					"dataset_source_id": config.StringVariable("2"),
					"dataset_query":     config.StringVariable("SELECT 1"),
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("censusworkspace_sql_dataset.test", plancheck.ResourceActionReplace),
						plancheck.ExpectUnknownValue("censusworkspace_sql_dataset.test", tfjsonpath.New("id")),
						plancheck.ExpectKnownValue("censusworkspace_sql_dataset.test", tfjsonpath.New("name"), knownvalue.StringExact("Test SQL Dataset (replaced)")),
						plancheck.ExpectKnownValue("censusworkspace_sql_dataset.test", tfjsonpath.New("query"), knownvalue.StringExact("SELECT 1")),
					},
					PostApplyPreRefresh: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue("censusworkspace_sql_dataset.test", tfjsonpath.New("id"), knownvalue.NotNull()),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue("censusworkspace_sql_dataset.test", tfjsonpath.New("id"), knownvalue.NotNull()),
					},
				},
			},
		},
	})
}
