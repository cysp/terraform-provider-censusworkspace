package provider_test

import (
	"net/http/httptest"
	"os"
	"testing"

	"github.com/cysp/terraform-provider-censusworkspace/internal/provider"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func CensusProviderMockedResourceTest(t *testing.T, testserver *httptest.Server, testcase resource.TestCase) {
	t.Helper()

	censusProviderMockableResourceTest(t, testserver, true, testcase)
}

func CensusProviderMockableResourceTest(t *testing.T, testserver *httptest.Server, testcase resource.TestCase) {
	t.Helper()

	censusProviderMockableResourceTest(t, testserver, false, testcase)
}

func censusProviderMockableResourceTest(t *testing.T, testserver *httptest.Server, alwaysMock bool, testcase resource.TestCase) {
	t.Helper()

	switch {
	case alwaysMock || os.Getenv("TF_ACC_MOCKED") != "":
		if testcase.ProtoV6ProviderFactories != nil {
			t.Fatal("tc.ProtoV6ProviderFactories must be nil")
		}

		testcase.ProtoV6ProviderFactories = makeTestAccProtoV6ProviderFactories(CensusProviderOptionsWithHTTPTestServer(testserver)...)
		resource.Test(t, testcase)

	default:
		if testcase.ProtoV6ProviderFactories == nil {
			testcase.ProtoV6ProviderFactories = testAccProtoV6ProviderFactories
		}

		resource.Test(t, testcase)
	}
}

func CensusProviderOptionsWithHTTPTestServer(testserver *httptest.Server) []provider.Option {
	if testserver == nil {
		return nil
	}

	return []provider.Option{
		provider.WithCensusURL(testserver.URL),
		provider.WithHTTPClient(testserver.Client()),
		provider.WithAccessToken("12345"),
	}
}
