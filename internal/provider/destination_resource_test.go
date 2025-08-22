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
func TestAccDestinationResourceImport(t *testing.T) {
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
				ResourceName:    "censusworkspace_destination.test",
				ImportState:     true,
				ImportStateId:   testDestinationIDString,
			},
		},
	})
}

//nolint:paralleltest
func TestAccDestinationResourceImportNotFound(t *testing.T) {
	server, err := cmt.NewCensusManagementServer()
	require.NoError(t, err)

	ProviderMockedResourceTest(t, server, resource.TestCase{
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestNameDirectory(),
				ResourceName:    "censusworkspace_destination.test",
				ImportState:     true,
				ImportStateId:   "99999",
				ExpectError:     regexp.MustCompile(`Cannot import non-existent remote object`),
			},
		},
	})
}

//nolint:paralleltest
func TestAccDestinationResourceCreateUpdateDelete(t *testing.T) {
	server, err := cmt.NewCensusManagementServer()
	require.NoError(t, err)

	ProviderMockedResourceTest(t, server, resource.TestCase{
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"destination_type": config.StringVariable("custom_api"),
					"destination_name": config.StringVariable("Test Destination"),
					"destination_credentials": config.MapVariable(map[string]config.Variable{
						"webhook_url": config.StringVariable("https://example.org/census-destination"),
					}),
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("censusworkspace_destination.test", plancheck.ResourceActionCreate),
						plancheck.ExpectUnknownValue("censusworkspace_destination.test", tfjsonpath.New("id")),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue("censusworkspace_destination.test", tfjsonpath.New("last_test_succeeded"), knownvalue.Null()),
					},
				},
			},
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
					"destination_type": config.StringVariable("custom_api"),
					"destination_name": config.StringVariable("Test Destination"),
					"destination_credentials": config.MapVariable(map[string]config.Variable{
						"webhook_url": config.StringVariable("https://example.org/census-destination"),
					}),
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
						plancheck.ExpectResourceAction("censusworkspace_destination.test", plancheck.ResourceActionNoop),
						plancheck.ExpectKnownValue("censusworkspace_destination.test", tfjsonpath.New("id"), knownvalue.NotNull()),
						plancheck.ExpectKnownValue("censusworkspace_destination.test", tfjsonpath.New("name"), knownvalue.NotNull()),
						plancheck.ExpectKnownValue("censusworkspace_destination.test", tfjsonpath.New("connection_details"), knownvalue.NotNull()),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue("censusworkspace_destination.test", tfjsonpath.New("last_test_succeeded"), knownvalue.Bool(false)),
					},
				},
			},
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"destination_type": config.StringVariable("custom_api"),
					"destination_name": config.StringVariable("Test Destination (updated)"),
					"destination_credentials": config.MapVariable(map[string]config.Variable{
						"webhook_url": config.StringVariable("https://example.org/census-destination"),
					}),
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("censusworkspace_destination.test", plancheck.ResourceActionUpdate),
						plancheck.ExpectKnownValue("censusworkspace_destination.test", tfjsonpath.New("id"), knownvalue.NotNull()),
						plancheck.ExpectKnownValue("censusworkspace_destination.test", tfjsonpath.New("name"), knownvalue.StringExact("Test Destination (updated)")),
						plancheck.ExpectUnknownValue("censusworkspace_destination.test", tfjsonpath.New("connection_details")),
					},
				},
			},
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"destination_type": config.StringVariable("big_query"),
					"destination_name": config.StringVariable("Test Destination (replaced)"),
					"destination_credentials": config.MapVariable(map[string]config.Variable{
						"project_id": config.StringVariable("project-id"),
					}),
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("censusworkspace_destination.test", plancheck.ResourceActionReplace),
						plancheck.ExpectUnknownValue("censusworkspace_destination.test", tfjsonpath.New("id")),
						plancheck.ExpectKnownValue("censusworkspace_destination.test", tfjsonpath.New("name"), knownvalue.StringExact("Test Destination (replaced)")),
						plancheck.ExpectUnknownValue("censusworkspace_destination.test", tfjsonpath.New("connection_details")),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue("censusworkspace_destination.test", tfjsonpath.New("last_test_succeeded"), knownvalue.Null()),
					},
				},
			},
		},
	})
}
