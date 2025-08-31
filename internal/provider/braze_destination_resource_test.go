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
func TestAccBrazeDestinationResourceImport(t *testing.T) {
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
		Type: "braze",
	}

	ProviderMockedResourceTest(t, server, resource.TestCase{
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: configVariables,
				ResourceName:    "censusworkspace_braze_destination.test",
				ImportState:     true,
				ImportStateId:   testDestinationIDString,
			},
		},
	})
}

//nolint:paralleltest
func TestAccBrazeDestinationResourceImportNotFound(t *testing.T) {
	server, err := cmt.NewCensusManagementServer()
	require.NoError(t, err)

	ProviderMockedResourceTest(t, server, resource.TestCase{
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestNameDirectory(),
				ResourceName:    "censusworkspace_braze_destination.test",
				ImportState:     true,
				ImportStateId:   "99999",
				ExpectError:     regexp.MustCompile(`Cannot import non-existent remote object`),
			},
		},
	})
}

//nolint:dupl,paralleltest
func TestAccBrazeDestinationResourceCreateUpdateDelete(t *testing.T) {
	server, err := cmt.NewCensusManagementServer()
	require.NoError(t, err)

	ProviderMockedResourceTest(t, server, resource.TestCase{
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"destination_type": config.StringVariable("braze"),
					"destination_name": config.StringVariable("Test Destination"),
					"destination_credentials": config.ObjectVariable(map[string]config.Variable{
						"instance_url": config.StringVariable("instance-url"),
						"api_key":      config.StringVariable("api-key"),
					}),
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("censusworkspace_braze_destination.test", plancheck.ResourceActionCreate),
						plancheck.ExpectUnknownValue("censusworkspace_braze_destination.test", tfjsonpath.New("id")),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue("censusworkspace_braze_destination.test", tfjsonpath.New("last_test_succeeded"), knownvalue.Null()),
					},
				},
			},
			//nolint:dupl
			{
				PreConfig: func() {
					destination := server.Handler().Destinations["1"]
					if destination != nil {
						destination.LastTestSucceeded.SetTo(false)
						destination.LastTestedAt.SetTo(time.Unix(0, 0))
					}
				},
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"destination_type": config.StringVariable("braze"),
					"destination_name": config.StringVariable("Test Destination"),
					"destination_credentials": config.ObjectVariable(map[string]config.Variable{
						"instance_url": config.StringVariable("instance-url"),
						"api_key":      config.StringVariable("api-key"),
					}),
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
						plancheck.ExpectResourceAction("censusworkspace_braze_destination.test", plancheck.ResourceActionNoop),
						plancheck.ExpectKnownValue("censusworkspace_braze_destination.test", tfjsonpath.New("id"), knownvalue.NotNull()),
						plancheck.ExpectKnownValue("censusworkspace_braze_destination.test", tfjsonpath.New("name"), knownvalue.NotNull()),
						plancheck.ExpectKnownValue("censusworkspace_braze_destination.test", tfjsonpath.New("connection_details"), knownvalue.NotNull()),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue("censusworkspace_braze_destination.test", tfjsonpath.New("last_test_succeeded"), knownvalue.Bool(false)),
					},
				},
			},
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"destination_type": config.StringVariable("braze"),
					"destination_name": config.StringVariable("Test Destination (updated)"),
					"destination_credentials": config.ObjectVariable(map[string]config.Variable{
						"instance_url": config.StringVariable("instance-url"),
						"api_key":      config.StringVariable("api-key"),
						"client_key":   config.StringVariable("client-key"),
					}),
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("censusworkspace_braze_destination.test", plancheck.ResourceActionUpdate),
						plancheck.ExpectKnownValue("censusworkspace_braze_destination.test", tfjsonpath.New("id"), knownvalue.NotNull()),
						plancheck.ExpectKnownValue("censusworkspace_braze_destination.test", tfjsonpath.New("name"), knownvalue.StringExact("Test Destination (updated)")),
						plancheck.ExpectUnknownValue("censusworkspace_braze_destination.test", tfjsonpath.New("connection_details")),
					},
				},
			},
		},
	})
}

//nolint:paralleltest
func TestAccBrazeDestinationResourceMovedFromDestination(t *testing.T) {
	server, err := cmt.NewCensusManagementServer()
	require.NoError(t, err)

	configVariables := config.Variables{
		"destination_name": config.StringVariable("Test Destination"),
		"destination_credentials": config.ObjectVariable(map[string]config.Variable{
			"instance_url": config.StringVariable("instance-url"),
			"api_key":      config.StringVariable("api-key"),
		}),
	}

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
						plancheck.ExpectKnownValue("censusworkspace_destination.test", tfjsonpath.New("last_test_succeeded"), knownvalue.Null()),
					},
				},
			},
			{
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: configVariables,

				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
						plancheck.ExpectResourceAction("censusworkspace_braze_destination.test", plancheck.ResourceActionNoop),
						plancheck.ExpectKnownValue("censusworkspace_braze_destination.test", tfjsonpath.New("id"), knownvalue.NotNull()),
						plancheck.ExpectKnownValue("censusworkspace_braze_destination.test", tfjsonpath.New("name"), knownvalue.StringExact("Test Destination")),
						plancheck.ExpectKnownValue("censusworkspace_braze_destination.test", tfjsonpath.New("credentials").AtMapKey("instance_url"), knownvalue.StringExact("instance-url")),
						plancheck.ExpectKnownValue("censusworkspace_braze_destination.test", tfjsonpath.New("credentials").AtMapKey("api_key"), knownvalue.StringExact("api-key")),
						plancheck.ExpectKnownValue("censusworkspace_braze_destination.test", tfjsonpath.New("connection_details"), knownvalue.NotNull()),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue("censusworkspace_braze_destination.test", tfjsonpath.New("last_test_succeeded"), knownvalue.Null()),
					},
				},
			},
		},
	})
}
