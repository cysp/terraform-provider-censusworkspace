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
func TestAccBigQuerySourceResourceImport(t *testing.T) {
	server, err := cmt.NewCensusManagementServer()
	require.NoError(t, err)

	testSourceID := int64(12345)
	testSourceIDString := strconv.FormatInt(testSourceID, 10)

	configVariables := config.Variables{
		"source_id": config.StringVariable(testSourceIDString),
	}

	server.Handler().Sources[testSourceIDString] = &cm.SourceData{
		ID:   testSourceID,
		Name: "Test_Source",
	}

	ProviderMockedResourceTest(t, server, resource.TestCase{
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: configVariables,
				ResourceName:    "censusworkspace_big_query_source.test",
				ImportState:     true,
				ImportStateId:   testSourceIDString,
			},
		},
	})
}

//nolint:paralleltest
func TestAccBigQuerySourceResourceImportNotFound(t *testing.T) {
	server, err := cmt.NewCensusManagementServer()
	require.NoError(t, err)

	ProviderMockedResourceTest(t, server, resource.TestCase{
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestNameDirectory(),
				ResourceName:    "censusworkspace_big_query_source.test",
				ImportState:     true,
				ImportStateId:   "99999",
				ExpectError:     regexp.MustCompile(`Cannot import non-existent remote object`),
			},
		},
	})
}

//nolint:paralleltest
func TestAccBigQuerySourceResourceCreateUpdateDelete(t *testing.T) {
	server, err := cmt.NewCensusManagementServer()
	require.NoError(t, err)

	ProviderMockedResourceTest(t, server, resource.TestCase{
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"source_name": config.StringVariable("bigquery_xxx_census_dev"),
					"source_credentials": config.ObjectVariable(map[string]config.Variable{
						"project_id": config.StringVariable("project-id"),
						"location":   config.StringVariable("US"),
					}),
					"source_warehouse_writeback_retention_in_days": config.IntegerVariable(7),
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("censusworkspace_big_query_source.test", plancheck.ResourceActionCreate),
						plancheck.ExpectUnknownValue("censusworkspace_big_query_source.test", tfjsonpath.New("id")),
						plancheck.ExpectKnownValue("censusworkspace_big_query_source.test", tfjsonpath.New("name"), knownvalue.StringExact("bigquery_xxx_census_dev")),
						plancheck.ExpectKnownValue("censusworkspace_big_query_source.test", tfjsonpath.New("warehouse_writeback_retention_in_days"), knownvalue.Int64Exact(7)),
					},
				},
			},
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"source_name": config.StringVariable("bigquery_xxx_census_dev"),
					"source_credentials": config.ObjectVariable(map[string]config.Variable{
						"project_id": config.StringVariable("project-id"),
						"location":   config.StringVariable("US"),
					}),
					"source_warehouse_writeback_retention_in_days": config.IntegerVariable(7),
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
						plancheck.ExpectResourceAction("censusworkspace_big_query_source.test", plancheck.ResourceActionNoop),
						plancheck.ExpectKnownValue("censusworkspace_big_query_source.test", tfjsonpath.New("id"), knownvalue.NotNull()),
						plancheck.ExpectKnownValue("censusworkspace_big_query_source.test", tfjsonpath.New("name"), knownvalue.StringExact("bigquery_xxx_census_dev")),
						plancheck.ExpectKnownValue("censusworkspace_big_query_source.test", tfjsonpath.New("warehouse_writeback_retention_in_days"), knownvalue.Int64Exact(7)),
						plancheck.ExpectKnownValue("censusworkspace_big_query_source.test", tfjsonpath.New("connection_details"), knownvalue.NotNull()),
					},
				},
			},
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"source_name": config.StringVariable("bigquery_xxx_census_prod"),
					"source_credentials": config.ObjectVariable(map[string]config.Variable{
						"project_id": config.StringVariable("project-id"),
						"location":   config.StringVariable("US"),
						"service_account_key": config.ObjectVariable(map[string]config.Variable{
							"type":           config.StringVariable("service_account"),
							"project_id":     config.StringVariable("project-id"),
							"private_key_id": config.StringVariable("private-key-id"),
							"private_key":    config.StringVariable("private-key"),
							"client_id":      config.StringVariable("client-id"),
							"client_email":   config.StringVariable("client-email"),
						}),
					}),
					"source_warehouse_writeback_retention_in_days": config.IntegerVariable(7),
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("censusworkspace_big_query_source.test", plancheck.ResourceActionUpdate),
						plancheck.ExpectKnownValue("censusworkspace_big_query_source.test", tfjsonpath.New("id"), knownvalue.NotNull()),
						plancheck.ExpectKnownValue("censusworkspace_big_query_source.test", tfjsonpath.New("name"), knownvalue.StringExact("bigquery_xxx_census_prod")),
						plancheck.ExpectKnownValue("censusworkspace_big_query_source.test", tfjsonpath.New("warehouse_writeback_retention_in_days"), knownvalue.Int64Exact(7)),
						plancheck.ExpectUnknownValue("censusworkspace_big_query_source.test", tfjsonpath.New("connection_details")),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue("censusworkspace_big_query_source.test", tfjsonpath.New("name"), knownvalue.StringExact("bigquery_xxx_census_prod")),
					},
				},
			},
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"source_name": config.StringVariable("bigquery_xxx_census_prod"),
					"source_credentials": config.ObjectVariable(map[string]config.Variable{
						"project_id": config.StringVariable("project-id"),
						"location":   config.StringVariable("US"),
						"service_account_key": config.ObjectVariable(map[string]config.Variable{
							"type":           config.StringVariable("service_account"),
							"project_id":     config.StringVariable("project-id"),
							"private_key_id": config.StringVariable("private-key-id"),
							"private_key":    config.StringVariable("private-key"),
							"client_id":      config.StringVariable("client-id"),
							"client_email":   config.StringVariable("client-email"),
						}),
					}),
					"source_warehouse_writeback_retention_in_days": config.IntegerVariable(7),
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
						plancheck.ExpectResourceAction("censusworkspace_big_query_source.test", plancheck.ResourceActionNoop),
						plancheck.ExpectKnownValue("censusworkspace_big_query_source.test", tfjsonpath.New("name"), knownvalue.StringExact("bigquery_xxx_census_prod")),
					},
				},
			},
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"source_name": config.StringVariable("bigquery_xxx_census_prod"),
					"source_credentials": config.ObjectVariable(map[string]config.Variable{
						"project_id": config.StringVariable("project-id"),
						"location":   config.StringVariable("US"),
					}),
					"source_warehouse_writeback_retention_in_days": config.IntegerVariable(14),
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("censusworkspace_big_query_source.test", plancheck.ResourceActionUpdate),
						plancheck.ExpectKnownValue("censusworkspace_big_query_source.test", tfjsonpath.New("warehouse_writeback_retention_in_days"), knownvalue.Int64Exact(14)),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue("censusworkspace_big_query_source.test", tfjsonpath.New("warehouse_writeback_retention_in_days"), knownvalue.Int64Exact(14)),
					},
				},
			},
		},
	})
}

//nolint:paralleltest
func TestAccBigQuerySourceResourceRejectsNonCanonicalName(t *testing.T) {
	server, err := cmt.NewCensusManagementServer()
	require.NoError(t, err)

	const canonicalConfig = `
resource "censusworkspace_big_query_source" "test" {
  name = "bigquery_xxx_census_dev"

  credentials = {
    project_id = "project-id"
    location   = "US"
  }
}
`

	const nonCanonicalConfig = `
resource "censusworkspace_big_query_source" "test" {
  name = "bigquery - xxx-census-dev"

  credentials = {
    project_id = "project-id"
    location   = "US"
  }
}
`

	ProviderMockedResourceTest(t, server, resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: canonicalConfig,
			},
			{
				Config:      nonCanonicalConfig,
				ExpectError: regexp.MustCompile(`Source names must use underscores`),
			},
			{
				Config: canonicalConfig,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
						plancheck.ExpectResourceAction("censusworkspace_big_query_source.test", plancheck.ResourceActionNoop),
					},
				},
			},
		},
	})
}

//nolint:paralleltest
func TestAccBigQuerySourceResourceMovedFromSource(t *testing.T) {
	server, err := cmt.NewCensusManagementServer()
	require.NoError(t, err)

	configVariables := config.Variables{
		"source_name": config.StringVariable("Test_Source"),
		"source_credentials": config.ObjectVariable(map[string]config.Variable{
			"project_id": config.StringVariable("project-id"),
			"location":   config.StringVariable("US"),
		}),
	}

	ProviderMockedResourceTest(t, server, resource.TestCase{
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: configVariables,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("censusworkspace_source.test", plancheck.ResourceActionCreate),
						plancheck.ExpectUnknownValue("censusworkspace_source.test", tfjsonpath.New("id")),
						plancheck.ExpectKnownValue("censusworkspace_source.test", tfjsonpath.New("name"), knownvalue.StringExact("Test_Source")),
					},
				},
			},
			{
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: configVariables,

				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
						plancheck.ExpectResourceAction("censusworkspace_big_query_source.test", plancheck.ResourceActionNoop),
						plancheck.ExpectKnownValue("censusworkspace_big_query_source.test", tfjsonpath.New("id"), knownvalue.NotNull()),
						plancheck.ExpectKnownValue("censusworkspace_big_query_source.test", tfjsonpath.New("name"), knownvalue.StringExact("Test_Source")),
						plancheck.ExpectKnownValue("censusworkspace_big_query_source.test", tfjsonpath.New("credentials").AtMapKey("project_id"), knownvalue.StringExact("project-id")),
						plancheck.ExpectKnownValue("censusworkspace_big_query_source.test", tfjsonpath.New("credentials").AtMapKey("location"), knownvalue.StringExact("US")),
						plancheck.ExpectKnownValue("censusworkspace_big_query_source.test", tfjsonpath.New("credentials").AtMapKey("service_account_key"), knownvalue.Null()),
						plancheck.ExpectKnownValue("censusworkspace_big_query_source.test", tfjsonpath.New("connection_details"), knownvalue.NotNull()),
					},
				},
			},
		},
	})
}

//nolint:paralleltest
func TestAccBigQuerySourceResourceMovedFromSourceWithServiceAccountKey(t *testing.T) {
	server, err := cmt.NewCensusManagementServer()
	require.NoError(t, err)

	configVariables := config.Variables{
		"source_name": config.StringVariable("Test_Source"),
		"source_credentials": config.ObjectVariable(map[string]config.Variable{
			"project_id": config.StringVariable("project-id"),
			"location":   config.StringVariable("US"),
			"service_account_key": config.ObjectVariable(map[string]config.Variable{
				"type":           config.StringVariable("service_account"),
				"project_id":     config.StringVariable("project-id"),
				"private_key_id": config.StringVariable("private-key-id"),
				"private_key":    config.StringVariable("private-key"),
				"client_id":      config.StringVariable("client-id"),
				"client_email":   config.StringVariable("client-email"),
			}),
		}),
	}

	ProviderMockedResourceTest(t, server, resource.TestCase{
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: configVariables,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("censusworkspace_source.test", plancheck.ResourceActionCreate),
						plancheck.ExpectUnknownValue("censusworkspace_source.test", tfjsonpath.New("id")),
						plancheck.ExpectKnownValue("censusworkspace_source.test", tfjsonpath.New("name"), knownvalue.StringExact("Test_Source")),
					},
				},
			},
			{
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: configVariables,

				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
						plancheck.ExpectResourceAction("censusworkspace_big_query_source.test", plancheck.ResourceActionNoop),
						plancheck.ExpectKnownValue("censusworkspace_big_query_source.test", tfjsonpath.New("id"), knownvalue.NotNull()),
						plancheck.ExpectKnownValue("censusworkspace_big_query_source.test", tfjsonpath.New("name"), knownvalue.StringExact("Test_Source")),
						plancheck.ExpectKnownValue("censusworkspace_big_query_source.test", tfjsonpath.New("credentials").AtMapKey("project_id"), knownvalue.StringExact("project-id")),
						plancheck.ExpectKnownValue("censusworkspace_big_query_source.test", tfjsonpath.New("credentials").AtMapKey("location"), knownvalue.StringExact("US")),
						plancheck.ExpectKnownValue("censusworkspace_big_query_source.test", tfjsonpath.New("credentials").AtMapKey("service_account_key").AtMapKey("private_key_id"), knownvalue.StringExact("private-key-id")),
						plancheck.ExpectKnownValue("censusworkspace_big_query_source.test", tfjsonpath.New("credentials").AtMapKey("service_account_key").AtMapKey("client_email"), knownvalue.StringExact("client-email")),
						plancheck.ExpectKnownValue("censusworkspace_big_query_source.test", tfjsonpath.New("connection_details"), knownvalue.NotNull()),
					},
				},
			},
		},
	})
}
