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
func TestAccCustomAPIDestinationResourceImport(t *testing.T) {
	server, err := cmt.NewCensusManagementServer()
	require.NoError(t, err)

	testDestinationID := int64(12345)
	testDestinationIDString := strconv.FormatInt(testDestinationID, 10)

	configVariables := config.Variables{
		"destination_id": config.StringVariable(testDestinationIDString),
	}

	server.Handler().Destinations[testDestinationIDString] = &cm.DestinationData{
		ID:   testDestinationID,
		Name: "Test Destination",
		Type: "custom_api",
	}

	ProviderMockedResourceTest(t, server, resource.TestCase{
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: configVariables,
				ResourceName:    "censusworkspace_custom_api_destination.test",
				ImportState:     true,
				ImportStateId:   testDestinationIDString,
			},
		},
	})
}

//nolint:paralleltest
func TestAccCustomAPIDestinationResourceImportNotFound(t *testing.T) {
	server, err := cmt.NewCensusManagementServer()
	require.NoError(t, err)

	ProviderMockedResourceTest(t, server, resource.TestCase{
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestNameDirectory(),
				ResourceName:    "censusworkspace_custom_api_destination.test",
				ImportState:     true,
				ImportStateId:   "99999",
				ExpectError:     regexp.MustCompile(`Cannot import non-existent remote object`),
			},
		},
	})
}

//nolint:paralleltest
func TestAccCustomAPIDestinationResourceCreateUpdateDelete(t *testing.T) {
	server, err := cmt.NewCensusManagementServer()
	require.NoError(t, err)

	ProviderMockedResourceTest(t, server, resource.TestCase{
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"destination_type": config.StringVariable("custom_api"),
					"destination_name": config.StringVariable("Test Destination"),
					"destination_credentials": config.ObjectVariable(map[string]config.Variable{
						"api_version": config.IntegerVariable(1),
						"webhook_url": config.StringVariable("https://example.org/census-destination"),
					}),
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("censusworkspace_custom_api_destination.test", plancheck.ResourceActionCreate),
						plancheck.ExpectUnknownValue("censusworkspace_custom_api_destination.test", tfjsonpath.New("id")),
					},
				},
			},
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"destination_type": config.StringVariable("custom_api"),
					"destination_name": config.StringVariable("Test Destination"),
					"destination_credentials": config.ObjectVariable(map[string]config.Variable{
						"api_version": config.IntegerVariable(1),
						"webhook_url": config.StringVariable("https://example.org/census-destination"),
					}),
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
						plancheck.ExpectResourceAction("censusworkspace_custom_api_destination.test", plancheck.ResourceActionNoop),
						plancheck.ExpectKnownValue("censusworkspace_custom_api_destination.test", tfjsonpath.New("id"), knownvalue.NotNull()),
						plancheck.ExpectKnownValue("censusworkspace_custom_api_destination.test", tfjsonpath.New("name"), knownvalue.NotNull()),
						plancheck.ExpectKnownValue("censusworkspace_custom_api_destination.test", tfjsonpath.New("connection_details"), knownvalue.NotNull()),
					},
				},
			},
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"destination_type": config.StringVariable("custom_api"),
					"destination_name": config.StringVariable("Test Destination (updated)"),
					"destination_credentials": config.ObjectVariable(map[string]config.Variable{
						"api_version": config.IntegerVariable(1),
						"webhook_url": config.StringVariable("https://example.org/census-destination"),
						"custom_headers": config.MapVariable(map[string]config.Variable{
							"x-client-id": config.ObjectVariable(map[string]config.Variable{
								"value":     config.StringVariable("client-123"),
								"is_secret": config.BoolVariable(false),
							}),
						}),
					}),
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("censusworkspace_custom_api_destination.test", plancheck.ResourceActionUpdate),
						plancheck.ExpectKnownValue("censusworkspace_custom_api_destination.test", tfjsonpath.New("id"), knownvalue.NotNull()),
						plancheck.ExpectKnownValue("censusworkspace_custom_api_destination.test", tfjsonpath.New("name"), knownvalue.StringExact("Test Destination (updated)")),
						plancheck.ExpectKnownValue("censusworkspace_custom_api_destination.test", tfjsonpath.New("credentials").AtMapKey("custom_headers").AtMapKey("x-client-id").AtMapKey("value"), knownvalue.StringExact("client-123")),
						plancheck.ExpectKnownValue("censusworkspace_custom_api_destination.test", tfjsonpath.New("credentials").AtMapKey("custom_headers").AtMapKey("x-client-id").AtMapKey("is_secret"), knownvalue.Bool(false)),
						plancheck.ExpectUnknownValue("censusworkspace_custom_api_destination.test", tfjsonpath.New("connection_details")),
					},
				},
			},
		},
	})
}

//nolint:paralleltest
func TestAccCustomAPIDestinationResourceWriteOnlyHeaders(t *testing.T) {
	server, err := cmt.NewCensusManagementServer()
	require.NoError(t, err)

	ProviderMockedResourceTest(t, server, resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: `
resource "censusworkspace_custom_api_destination" "test" {
  name = "Test Destination"

  credentials = {
    api_version = 1
    webhook_url = "https://example.org/census-destination"
    custom_headers = {
      "x-client-id" = {
        value     = "client-123"
        is_secret = false
      }
      "x-\"secret\"\\key" = {
        value     = "secret"
        is_secret = true
      }
    }
  }
}
`,
			},
			{
				Config: `
resource "censusworkspace_custom_api_destination" "test" {
  name = "Test Destination"

  credentials = {
    api_version = 1
    webhook_url = "https://example.org/census-destination"
    custom_headers = {
      "x-client-id" = {
        value     = "client-123"
        is_secret = false
      }
      "x-\"secret\"\\key" = {
        value_wo  = "secret"
        is_secret = true
      }
    }
  }
}
`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("censusworkspace_custom_api_destination.test", plancheck.ResourceActionUpdate),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue("censusworkspace_custom_api_destination.test", tfjsonpath.New("credentials").AtMapKey("custom_headers").AtMapKey(`x-"secret"\key`).AtMapKey("value"), knownvalue.Null()),
						plancheck.ExpectKnownValue("censusworkspace_custom_api_destination.test", tfjsonpath.New("credentials").AtMapKey("custom_headers").AtMapKey("x-client-id").AtMapKey("value"), knownvalue.StringExact("client-123")),
					},
				},
			},
			{
				Config: `
resource "censusworkspace_custom_api_destination" "test" {
  name = "Test Destination"

  credentials = {
    api_version = 1
    webhook_url = "https://example.org/census-destination"
    custom_headers = {
      "x-client-id" = {
        value     = "client-123"
        is_secret = false
      }
      "x-\"secret\"\\key" = {
        value_wo  = "rotated-secret"
        is_secret = true
      }
    }
  }
}
`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("censusworkspace_custom_api_destination.test", plancheck.ResourceActionUpdate),
						plancheck.ExpectUnknownValue("censusworkspace_custom_api_destination.test", tfjsonpath.New("connection_details")),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue("censusworkspace_custom_api_destination.test", tfjsonpath.New("credentials").AtMapKey("custom_headers").AtMapKey(`x-"secret"\key`).AtMapKey("value"), knownvalue.Null()),
						plancheck.ExpectKnownValue("censusworkspace_custom_api_destination.test", tfjsonpath.New("credentials").AtMapKey("custom_headers").AtMapKey("x-client-id").AtMapKey("value"), knownvalue.StringExact("client-123")),
					},
				},
			},
			{
				Config: `
resource "censusworkspace_custom_api_destination" "test" {
  name = "Test Destination"

  credentials = {
    api_version = 1
    webhook_url = "https://example.org/census-destination"
    custom_headers = {
      "x-client-id" = {
        value     = "client-123"
        is_secret = false
      }
      "x-\"secret\"\\key" = {
        value_wo  = "rotated-secret"
        is_secret = true
      }
    }
  }
}
`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
						plancheck.ExpectResourceAction("censusworkspace_custom_api_destination.test", plancheck.ResourceActionNoop),
					},
				},
			},
			{
				Config: `
resource "censusworkspace_custom_api_destination" "test" {
  name = "Test Destination (updated)"

  credentials = {
    api_version = 1
    webhook_url = "https://example.org/census-destination"
    custom_headers = {
      "x-client-id" = {
        value     = "client-123"
        is_secret = false
      }
      "x-\"secret\"\\key" = {
        value_wo  = "rotated-secret"
        is_secret = true
      }
    }
  }
}
`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("censusworkspace_custom_api_destination.test", plancheck.ResourceActionUpdate),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue("censusworkspace_custom_api_destination.test", tfjsonpath.New("name"), knownvalue.StringExact("Test Destination (updated)")),
						plancheck.ExpectKnownValue("censusworkspace_custom_api_destination.test", tfjsonpath.New("credentials").AtMapKey("custom_headers").AtMapKey(`x-"secret"\key`).AtMapKey("value"), knownvalue.Null()),
						plancheck.ExpectKnownValue("censusworkspace_custom_api_destination.test", tfjsonpath.New("credentials").AtMapKey("custom_headers").AtMapKey("x-client-id").AtMapKey("value"), knownvalue.StringExact("client-123")),
					},
				},
			},
		},
	})
}

//nolint:paralleltest
func TestAccCustomAPIDestinationResourceMissingHeaderValue(t *testing.T) {
	server, err := cmt.NewCensusManagementServer()
	require.NoError(t, err)

	ProviderMockedResourceTest(t, server, resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: `
resource "censusworkspace_custom_api_destination" "test" {
  name = "Test Destination"

  credentials = {
    api_version = 1
    webhook_url = "https://example.org/census-destination"
    custom_headers = {
      "x-client-secret" = {
        is_secret = true
      }
    }
  }
}
`,
				ExpectError: regexp.MustCompile(`One of credentials\.custom_headers\["x-client-secret"\]\.value or\s+credentials\.custom_headers\["x-client-secret"\]\.value_wo must be configured`),
			},
		},
	})
}

//nolint:dupl,paralleltest
func TestAccCustomAPIDestinationResourceMovedFromDestination(t *testing.T) {
	server, err := cmt.NewCensusManagementServer()
	require.NoError(t, err)

	configVariables := config.Variables{
		"destination_name": config.StringVariable("Test Destination"),
		"destination_credentials": config.ObjectVariable(map[string]config.Variable{
			"api_version": config.IntegerVariable(1),
			"webhook_url": config.StringVariable("https://example.org/census-destination"),
		}),
	}

	//nolint:dupl
	ProviderMockedResourceTest(t, server, resource.TestCase{
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: configVariables,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("censusworkspace_destination.test", plancheck.ResourceActionCreate),
						plancheck.ExpectUnknownValue("censusworkspace_destination.test", tfjsonpath.New("id")),
						plancheck.ExpectKnownValue("censusworkspace_destination.test", tfjsonpath.New("name"), knownvalue.StringExact("Test Destination")),
						plancheck.ExpectUnknownValue("censusworkspace_destination.test", tfjsonpath.New("connection_details")),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue("censusworkspace_destination.test", tfjsonpath.New("connection_details"), knownvalue.NotNull()),
					},
				},
			},
			{
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: configVariables,

				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
						plancheck.ExpectResourceAction("censusworkspace_custom_api_destination.test", plancheck.ResourceActionNoop),
						plancheck.ExpectKnownValue("censusworkspace_custom_api_destination.test", tfjsonpath.New("id"), knownvalue.NotNull()),
						plancheck.ExpectKnownValue("censusworkspace_custom_api_destination.test", tfjsonpath.New("name"), knownvalue.StringExact("Test Destination")),
						plancheck.ExpectKnownValue("censusworkspace_custom_api_destination.test", tfjsonpath.New("credentials").AtMapKey("api_version"), knownvalue.Int64Exact(1)),
						plancheck.ExpectKnownValue("censusworkspace_custom_api_destination.test", tfjsonpath.New("credentials").AtMapKey("webhook_url"), knownvalue.StringExact("https://example.org/census-destination")),
						plancheck.ExpectKnownValue("censusworkspace_custom_api_destination.test", tfjsonpath.New("credentials").AtMapKey("custom_headers"), knownvalue.Null()),
						plancheck.ExpectKnownValue("censusworkspace_custom_api_destination.test", tfjsonpath.New("connection_details"), knownvalue.NotNull()),
					},
				},
			},
		},
	})
}

//nolint:paralleltest
func TestAccCustomAPIDestinationResourceMovedFromDestinationWithCustomHeaders(t *testing.T) {
	server, err := cmt.NewCensusManagementServer()
	require.NoError(t, err)

	configVariables := config.Variables{
		"destination_name": config.StringVariable("Test Destination"),
		"destination_credentials": config.ObjectVariable(map[string]config.Variable{
			"api_version": config.IntegerVariable(1),
			"webhook_url": config.StringVariable("https://example.org/census-destination"),
			"custom_headers": config.MapVariable(map[string]config.Variable{
				"x-client-id": config.ObjectVariable(map[string]config.Variable{
					"value":     config.StringVariable("client-123"),
					"is_secret": config.BoolVariable(false),
				}),
			}),
		}),
	}

	//nolint:dupl
	ProviderMockedResourceTest(t, server, resource.TestCase{
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: configVariables,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("censusworkspace_destination.test", plancheck.ResourceActionCreate),
						plancheck.ExpectUnknownValue("censusworkspace_destination.test", tfjsonpath.New("id")),
						plancheck.ExpectKnownValue("censusworkspace_destination.test", tfjsonpath.New("name"), knownvalue.StringExact("Test Destination")),
						plancheck.ExpectUnknownValue("censusworkspace_destination.test", tfjsonpath.New("connection_details")),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue("censusworkspace_destination.test", tfjsonpath.New("connection_details"), knownvalue.NotNull()),
					},
				},
			},
			{
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: configVariables,

				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
						plancheck.ExpectResourceAction("censusworkspace_custom_api_destination.test", plancheck.ResourceActionNoop),
						plancheck.ExpectKnownValue("censusworkspace_custom_api_destination.test", tfjsonpath.New("id"), knownvalue.NotNull()),
						plancheck.ExpectKnownValue("censusworkspace_custom_api_destination.test", tfjsonpath.New("name"), knownvalue.StringExact("Test Destination")),
						plancheck.ExpectKnownValue("censusworkspace_custom_api_destination.test", tfjsonpath.New("credentials").AtMapKey("api_version"), knownvalue.Int64Exact(1)),
						plancheck.ExpectKnownValue("censusworkspace_custom_api_destination.test", tfjsonpath.New("credentials").AtMapKey("webhook_url"), knownvalue.StringExact("https://example.org/census-destination")),
						plancheck.ExpectKnownValue("censusworkspace_custom_api_destination.test", tfjsonpath.New("credentials").AtMapKey("custom_headers"), knownvalue.NotNull()),
						plancheck.ExpectKnownValue("censusworkspace_custom_api_destination.test", tfjsonpath.New("connection_details"), knownvalue.NotNull()),
					},
				},
			},
		},
	})
}
