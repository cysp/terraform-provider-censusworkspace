package provider_test

import (
	"regexp"
	"strconv"
	"testing"
	"time"

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
		Name: "Test Source",
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
					"source_label": config.StringVariable("Test Source"),
					"source_credentials": config.ObjectVariable(map[string]config.Variable{
						"project_id": config.StringVariable("project-id"),
						"location":   config.StringVariable("US"),
					}),
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("censusworkspace_big_query_source.test", plancheck.ResourceActionCreate),
						plancheck.ExpectUnknownValue("censusworkspace_big_query_source.test", tfjsonpath.New("id")),
						plancheck.ExpectKnownValue("censusworkspace_big_query_source.test", tfjsonpath.New("label"), knownvalue.StringExact("Test Source")),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue("censusworkspace_big_query_source.test", tfjsonpath.New("last_test_succeeded"), knownvalue.Null()),
					},
				},
			},
			{
				PreConfig: func() {
					source := server.Handler().Sources["1"]
					if source != nil {
						source.LastTestSucceeded.SetTo(false)
						source.LastTestedAt.SetTo(time.Unix(0, 0))
					}
				},
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"source_label": config.StringVariable("Test Source"),
					"source_credentials": config.ObjectVariable(map[string]config.Variable{
						"project_id": config.StringVariable("project-id"),
						"location":   config.StringVariable("US"),
					}),
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
						plancheck.ExpectResourceAction("censusworkspace_big_query_source.test", plancheck.ResourceActionNoop),
						plancheck.ExpectKnownValue("censusworkspace_big_query_source.test", tfjsonpath.New("id"), knownvalue.NotNull()),
						plancheck.ExpectKnownValue("censusworkspace_big_query_source.test", tfjsonpath.New("label"), knownvalue.StringExact("Test Source")),
						plancheck.ExpectKnownValue("censusworkspace_big_query_source.test", tfjsonpath.New("name"), knownvalue.NotNull()),
						plancheck.ExpectKnownValue("censusworkspace_big_query_source.test", tfjsonpath.New("connection_details"), knownvalue.NotNull()),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue("censusworkspace_big_query_source.test", tfjsonpath.New("last_test_succeeded"), knownvalue.Bool(false)),
					},
				},
			},
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"source_label": config.StringVariable("Test Source (updated)"),
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
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("censusworkspace_big_query_source.test", plancheck.ResourceActionUpdate),
						plancheck.ExpectKnownValue("censusworkspace_big_query_source.test", tfjsonpath.New("id"), knownvalue.NotNull()),
						plancheck.ExpectKnownValue("censusworkspace_big_query_source.test", tfjsonpath.New("label"), knownvalue.StringExact("Test Source (updated)")),
						plancheck.ExpectUnknownValue("censusworkspace_big_query_source.test", tfjsonpath.New("name")),
						plancheck.ExpectUnknownValue("censusworkspace_big_query_source.test", tfjsonpath.New("connection_details")),
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
		"source_label": config.StringVariable("Test Source"),
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
						plancheck.ExpectKnownValue("censusworkspace_source.test", tfjsonpath.New("label"), knownvalue.StringExact("Test Source")),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue("censusworkspace_source.test", tfjsonpath.New("last_test_succeeded"), knownvalue.Null()),
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
						plancheck.ExpectKnownValue("censusworkspace_big_query_source.test", tfjsonpath.New("label"), knownvalue.StringExact("Test Source")),
						plancheck.ExpectKnownValue("censusworkspace_big_query_source.test", tfjsonpath.New("name"), knownvalue.NotNull()),
						plancheck.ExpectKnownValue("censusworkspace_big_query_source.test", tfjsonpath.New("credentials").AtMapKey("project_id"), knownvalue.StringExact("project-id")),
						plancheck.ExpectKnownValue("censusworkspace_big_query_source.test", tfjsonpath.New("credentials").AtMapKey("location"), knownvalue.StringExact("US")),
						plancheck.ExpectKnownValue("censusworkspace_big_query_source.test", tfjsonpath.New("credentials").AtMapKey("service_account_key"), knownvalue.Null()),
						plancheck.ExpectKnownValue("censusworkspace_big_query_source.test", tfjsonpath.New("connection_details"), knownvalue.NotNull()),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue("censusworkspace_big_query_source.test", tfjsonpath.New("last_test_succeeded"), knownvalue.Null()),
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
		"source_label": config.StringVariable("Test Source"),
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
						plancheck.ExpectKnownValue("censusworkspace_source.test", tfjsonpath.New("label"), knownvalue.StringExact("Test Source")),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue("censusworkspace_source.test", tfjsonpath.New("last_test_succeeded"), knownvalue.Null()),
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
						plancheck.ExpectKnownValue("censusworkspace_big_query_source.test", tfjsonpath.New("label"), knownvalue.StringExact("Test Source")),
						plancheck.ExpectKnownValue("censusworkspace_big_query_source.test", tfjsonpath.New("name"), knownvalue.NotNull()),
						plancheck.ExpectKnownValue("censusworkspace_big_query_source.test", tfjsonpath.New("credentials").AtMapKey("project_id"), knownvalue.StringExact("project-id")),
						plancheck.ExpectKnownValue("censusworkspace_big_query_source.test", tfjsonpath.New("credentials").AtMapKey("location"), knownvalue.StringExact("US")),
						plancheck.ExpectKnownValue("censusworkspace_big_query_source.test", tfjsonpath.New("credentials").AtMapKey("service_account_key").AtMapKey("private_key_id"), knownvalue.StringExact("private-key-id")),
						plancheck.ExpectKnownValue("censusworkspace_big_query_source.test", tfjsonpath.New("credentials").AtMapKey("service_account_key").AtMapKey("client_email"), knownvalue.StringExact("client-email")),
						plancheck.ExpectKnownValue("censusworkspace_big_query_source.test", tfjsonpath.New("connection_details"), knownvalue.NotNull()),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue("censusworkspace_big_query_source.test", tfjsonpath.New("last_test_succeeded"), knownvalue.Null()),
					},
				},
			},
		},
	})
}
