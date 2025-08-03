package provider_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/cysp/terraform-provider-censusworkspace/internal/provider"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func ProviderMockedResourceTest(t *testing.T, server http.Handler, testcase resource.TestCase) {
	t.Helper()

	providerMockableResourceTest(t, server, true, testcase)
}

func ProviderMockableResourceTest(t *testing.T, server http.Handler, testcase resource.TestCase) {
	t.Helper()

	providerMockableResourceTest(t, server, false, testcase)
}

func providerMockableResourceTest(t *testing.T, server http.Handler, alwaysMock bool, testcase resource.TestCase) {
	t.Helper()

	switch {
	case alwaysMock || os.Getenv("TF_ACC_MOCKED") != "":
		if testcase.ProtoV6ProviderFactories != nil {
			t.Fatal("tc.ProtoV6ProviderFactories must be nil")
		}

		var hts *httptest.Server
		if server != nil {
			hts = httptest.NewServer(server)
			defer hts.Close()
		}

		testcase.ProtoV6ProviderFactories = makeTestAccProtoV6ProviderFactories(ProviderOptionsWithHTTPTestServer(hts)...)
		resource.Test(t, testcase)

	default:
		if testcase.ProtoV6ProviderFactories == nil {
			testcase.ProtoV6ProviderFactories = testAccProtoV6ProviderFactories
		}

		resource.Test(t, testcase)
	}
}

func ProviderOptionsWithHTTPTestServer(testserver *httptest.Server) []provider.Option {
	if testserver == nil {
		return nil
	}

	return []provider.Option{
		provider.WithBaseURL(testserver.URL),
		provider.WithHTTPClient(testserver.Client()),
		provider.WithWorkspaceApiKey("12345"),
	}
}
