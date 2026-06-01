package provider_test

import (
	"net/http/httptest"
	"regexp"
	"strconv"
	"testing"
	"time"

	cm "github.com/cysp/terraform-provider-censusworkspace/internal/census-management-go"
	cmt "github.com/cysp/terraform-provider-censusworkspace/internal/census-management-go/testing"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/stretchr/testify/require"
)

//nolint:paralleltest
func TestAccSourceResourceImport(t *testing.T) {
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
				ResourceName:    "censusworkspace_source.test",
				ImportState:     true,
				ImportStateId:   testSourceIDString,
			},
		},
	})
}

//nolint:paralleltest
func TestAccSourceResourceImportNotFound(t *testing.T) {
	server, err := cmt.NewCensusManagementServer()
	require.NoError(t, err)

	ProviderMockedResourceTest(t, server, resource.TestCase{
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestNameDirectory(),
				ResourceName:    "censusworkspace_source.test",
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
	require.NoError(t, err)

	ProviderMockedResourceTest(t, server, resource.TestCase{
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"source_type": config.StringVariable("big_query"),
					"source_name": config.StringVariable("Test Source"),
					"source_credentials": config.MapVariable(map[string]config.Variable{
						"project_id": config.StringVariable("project-id"),
						"location":   config.StringVariable("US"),
					}),
					"source_warehouse_writeback_retention_in_days": config.IntegerVariable(7),
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("censusworkspace_source.test", plancheck.ResourceActionCreate),
						plancheck.ExpectUnknownValue("censusworkspace_source.test", tfjsonpath.New("id")),
						plancheck.ExpectKnownValue("censusworkspace_source.test", tfjsonpath.New("name"), knownvalue.StringExact("Test Source")),
						plancheck.ExpectKnownValue("censusworkspace_source.test", tfjsonpath.New("warehouse_writeback_retention_in_days"), knownvalue.Int64Exact(7)),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue("censusworkspace_source.test", tfjsonpath.New("last_test_succeeded"), knownvalue.Null()),
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
					"source_type": config.StringVariable("big_query"),
					"source_name": config.StringVariable("Test Source"),
					"source_credentials": config.MapVariable(map[string]config.Variable{
						"project_id": config.StringVariable("project-id"),
						"location":   config.StringVariable("US"),
					}),
					"source_warehouse_writeback_retention_in_days": config.IntegerVariable(7),
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
						plancheck.ExpectResourceAction("censusworkspace_source.test", plancheck.ResourceActionNoop),
						plancheck.ExpectKnownValue("censusworkspace_source.test", tfjsonpath.New("id"), knownvalue.NotNull()),
						plancheck.ExpectKnownValue("censusworkspace_source.test", tfjsonpath.New("name"), knownvalue.StringExact("Test Source")),
						plancheck.ExpectKnownValue("censusworkspace_source.test", tfjsonpath.New("warehouse_writeback_retention_in_days"), knownvalue.Int64Exact(7)),
						plancheck.ExpectKnownValue("censusworkspace_source.test", tfjsonpath.New("connection_details"), knownvalue.NotNull()),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue("censusworkspace_source.test", tfjsonpath.New("last_test_succeeded"), knownvalue.Bool(false)),
					},
				},
			},
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"source_type": config.StringVariable("big_query"),
					"source_name": config.StringVariable("Test Source (updated)"),
					"source_credentials": config.MapVariable(map[string]config.Variable{
						"project_id": config.StringVariable("project-id"),
						"location":   config.StringVariable("US"),
					}),
					"source_warehouse_writeback_retention_in_days": config.IntegerVariable(7),
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("censusworkspace_source.test", plancheck.ResourceActionUpdate),
						plancheck.ExpectKnownValue("censusworkspace_source.test", tfjsonpath.New("id"), knownvalue.NotNull()),
						plancheck.ExpectKnownValue("censusworkspace_source.test", tfjsonpath.New("name"), knownvalue.StringExact("Test Source (updated)")),
						plancheck.ExpectKnownValue("censusworkspace_source.test", tfjsonpath.New("warehouse_writeback_retention_in_days"), knownvalue.Int64Exact(7)),
						plancheck.ExpectUnknownValue("censusworkspace_source.test", tfjsonpath.New("connection_details")),
					},
				},
			},
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"source_type": config.StringVariable("big_query"),
					"source_name": config.StringVariable("Test Source (updated)"),
					"source_credentials": config.MapVariable(map[string]config.Variable{
						"project_id": config.StringVariable("project-id"),
						"location":   config.StringVariable("US"),
					}),
					"source_warehouse_writeback_retention_in_days": config.IntegerVariable(14),
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("censusworkspace_source.test", plancheck.ResourceActionUpdate),
						plancheck.ExpectKnownValue("censusworkspace_source.test", tfjsonpath.New("warehouse_writeback_retention_in_days"), knownvalue.Int64Exact(14)),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue("censusworkspace_source.test", tfjsonpath.New("warehouse_writeback_retention_in_days"), knownvalue.Int64Exact(14)),
					},
				},
			},
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"source_type": config.StringVariable("pub_sub"),
					"source_name": config.StringVariable("Test Source (replaced)"),
					"source_credentials": config.MapVariable(map[string]config.Variable{
						"project_id": config.StringVariable("project-id"),
					}),
					"source_warehouse_writeback_retention_in_days": config.IntegerVariable(14),
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("censusworkspace_source.test", plancheck.ResourceActionReplace),
						plancheck.ExpectUnknownValue("censusworkspace_source.test", tfjsonpath.New("id")),
						plancheck.ExpectKnownValue("censusworkspace_source.test", tfjsonpath.New("name"), knownvalue.StringExact("Test Source (replaced)")),
						plancheck.ExpectUnknownValue("censusworkspace_source.test", tfjsonpath.New("connection_details")),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue("censusworkspace_source.test", tfjsonpath.New("last_test_succeeded"), knownvalue.Null()),
					},
				},
			},
		},
	})
}

//nolint:paralleltest
func TestAccSourceResourceSyncEngineAppearing(t *testing.T) {
	server, err := cmt.NewCensusManagementServer()
	require.NoError(t, err)

	testSourceID := int64(12345)
	testSourceIDString := strconv.FormatInt(testSourceID, 10)

	configVariables := config.Variables{
		"source_id": config.StringVariable(testSourceIDString),
	}

	testSource := cm.SourceData{
		ID:    testSourceID,
		Name:  "Test Source",
		Type:  "big_query",
		Label: cm.NewOptNilString("Test Source"),
	}

	server.Handler().Sources[testSourceIDString] = &testSource

	ProviderMockedResourceTest(t, server, resource.TestCase{
		Steps: []resource.TestStep{
			{
				ConfigDirectory:    config.TestNameDirectory(),
				ConfigVariables:    configVariables,
				ResourceName:       "censusworkspace_source.test",
				ImportState:        true,
				ImportStateId:      testSourceIDString,
				ImportStatePersist: true,
			},
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: configVariables,
				ResourceName:    "censusworkspace_source.test",
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("censusworkspace_source.test", plancheck.ResourceActionNoop),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckNoResourceAttr("censusworkspace_source.test", "sync_engine"),
				),
			},
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: configVariables,
				PreConfig: func() {
					testSource.SyncEngine.SetTo("basic")
				},
				ResourceName: "censusworkspace_source.test",
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("censusworkspace_source.test", "sync_engine", "basic"),
				),
			},
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: configVariables,
				ResourceName:    "censusworkspace_source.test",
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("censusworkspace_source.test", plancheck.ResourceActionNoop),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("censusworkspace_source.test", "sync_engine", "basic"),
				),
			},
		},
	})
}

//nolint:paralleltest
func TestProtocol6SourceResourceReadAcceptsExistingLabelState(t *testing.T) {
	testProtocol6SourceResourceReadAcceptsExistingLabelState(t, "censusworkspace_source", map[string]tftypes.Value{
		"type":                tftypes.NewValue(tftypes.String, "big_query"),
		"credentials":         tftypes.NewValue(tftypes.String, `{"project_id":"project-id","location":"US"}`),
		"connection_details":  tftypes.NewValue(tftypes.String, nil),
		"last_tested_at":      tftypes.NewValue(tftypes.String, nil),
		"last_test_succeeded": tftypes.NewValue(tftypes.Bool, nil),
	})
}

//nolint:paralleltest
func TestProtocol6BigQuerySourceResourceReadAcceptsExistingLabelState(t *testing.T) {
	serviceAccountKeyType := tftypes.Object{AttributeTypes: map[string]tftypes.Type{
		"type":           tftypes.String,
		"project_id":     tftypes.String,
		"private_key_id": tftypes.String,
		"private_key":    tftypes.String,
		"client_id":      tftypes.String,
		"client_email":   tftypes.String,
	}}

	testProtocol6SourceResourceReadAcceptsExistingLabelState(t, "censusworkspace_big_query_source", map[string]tftypes.Value{
		"credentials": tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{
			"project_id":          tftypes.String,
			"location":            tftypes.String,
			"service_account_key": serviceAccountKeyType,
		}}, map[string]tftypes.Value{
			"project_id":          tftypes.NewValue(tftypes.String, "project-id"),
			"location":            tftypes.NewValue(tftypes.String, "US"),
			"service_account_key": tftypes.NewValue(serviceAccountKeyType, nil),
		}),
		"connection_details": tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{
			"project_id":      tftypes.String,
			"location":        tftypes.String,
			"service_account": tftypes.String,
		}}, nil),
		"last_tested_at":      tftypes.NewValue(tftypes.String, nil),
		"last_test_succeeded": tftypes.NewValue(tftypes.Bool, nil),
	})
}

func testProtocol6SourceResourceReadAcceptsExistingLabelState(
	t *testing.T,
	resourceTypeName string,
	resourceState map[string]tftypes.Value,
) {
	t.Helper()

	server, err := cmt.NewCensusManagementServer()
	require.NoError(t, err)

	testSourceID := int64(12345)
	testSourceIDString := strconv.FormatInt(testSourceID, 10)

	server.Handler().Sources[testSourceIDString] = &cm.SourceData{
		ID:    testSourceID,
		Name:  "Test Source",
		Type:  "big_query",
		Label: cm.NewOptNilString("Remote Label"),
	}

	hts := httptest.NewServer(server)
	defer hts.Close()

	providerServer, err := makeTestAccProtoV6ProviderFactories(ProviderOptionsWithHTTPTestServer(hts)...,
	)["censusworkspace"]()
	require.NoError(t, err)

	providerSchemaResp, err := providerServer.GetProviderSchema(t.Context(), &tfprotov6.GetProviderSchemaRequest{})
	require.NoError(t, err)
	require.Empty(t, providerSchemaResp.Diagnostics)

	providerConfig, err := providerConfigDynamicValue(map[string]any{})
	require.NoError(t, err)

	configureResp, err := providerServer.ConfigureProvider(t.Context(), &tfprotov6.ConfigureProviderRequest{
		Config: &providerConfig,
	})
	require.NoError(t, err)
	require.Empty(t, configureResp.Diagnostics)

	sourceSchema := providerSchemaResp.ResourceSchemas[resourceTypeName]
	sourceType := sourceSchema.ValueType()

	resourceState["id"] = tftypes.NewValue(tftypes.String, testSourceIDString)
	resourceState["name"] = tftypes.NewValue(tftypes.String, "Test Source")
	resourceState["label"] = tftypes.NewValue(tftypes.String, "Legacy Label")
	resourceState["sync_engine"] = tftypes.NewValue(tftypes.String, nil)
	resourceState["warehouse_writeback_retention_in_days"] = tftypes.NewValue(tftypes.Number, nil)

	currentState, err := tfprotov6.NewDynamicValue(sourceType, tftypes.NewValue(sourceType, resourceState))
	require.NoError(t, err)

	readResp, err := providerServer.ReadResource(t.Context(), &tfprotov6.ReadResourceRequest{
		TypeName:     resourceTypeName,
		CurrentState: &currentState,
	})
	require.NoError(t, err)
	require.Empty(t, readResp.Diagnostics)
	require.NotNil(t, readResp.NewState)

	newState, err := readResp.NewState.Unmarshal(sourceType)
	require.NoError(t, err)

	var newStateAttrs map[string]tftypes.Value
	require.NoError(t, newState.As(&newStateAttrs))

	var label string
	require.NoError(t, newStateAttrs["label"].As(&label))
	require.Equal(t, "Remote Label", label)
}
