package testing_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	cm "github.com/cysp/terraform-provider-censusworkspace/internal/census-management-go"
	cmt "github.com/cysp/terraform-provider-censusworkspace/internal/census-management-go/testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCensusManagementServerGetApiV1(t *testing.T) {
	t.Parallel()

	server, serverErr := cmt.NewCensusManagementServer()
	require.NoError(t, serverErr, "failed to create server")

	ts := httptest.NewServer(server)
	defer ts.Close()

	client, clientErr := cm.NewClient(ts.URL, cm.NewWorkspaceAPIKeySecuritySource("test-api-key"))
	require.NoError(t, clientErr, "failed to create client")

	response, responseErr := client.GetApiV1(t.Context())

	assert.Nil(t, response, "response should be nil")

	var responseStatusMessage *cm.StatusResponseStatusCode
	if assert.ErrorAs(t, responseErr, &responseStatusMessage) {
		assert.Equal(t, http.StatusNotFound, responseStatusMessage.StatusCode, "expected status code to be 404")
	}
}
