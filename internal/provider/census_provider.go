package provider

import (
	"context"
	"net/http"
	"os"
	"time"

	cm "github.com/cysp/terraform-provider-censusworkspace/internal/census-management-go"
	"github.com/cysp/terraform-provider-censusworkspace/internal/provider/util"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func Factory(version string, options ...Option) func() provider.Provider {
	return func() provider.Provider {
		return New(version, options...)
	}
}

func New(version string, options ...Option) *CensusProvider {
	provider := CensusProvider{
		version: version,
	}

	for _, option := range options {
		option(&provider)
	}

	return &provider
}

type CensusProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string

	baseURL    string
	httpClient *http.Client
	apiKey     string
}

type CensusProviderModel struct {
	URL    types.String `tfsdk:"url"`
	ApiKey types.String `tfsdk:"api_key"`
}

var _ provider.Provider = (*CensusProvider)(nil)

type Option func(*CensusProvider)

func WithCensusURL(url string) Option {
	return func(p *CensusProvider) {
		p.baseURL = url
	}
}

func WithHTTPClient(httpClient *http.Client) Option {
	return func(p *CensusProvider) {
		p.httpClient = httpClient
	}
}

func WithApiKey(apiKey string) Option {
	return func(p *CensusProvider) {
		p.apiKey = apiKey
	}
}

func (p *CensusProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage Census workspace.",
		Attributes: map[string]schema.Attribute{
			"url": schema.StringAttribute{
				Optional: true,
			},
			"api_key": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
		},
	}
}

func (p *CensusProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data CensusProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var censusBaseURL string
	if !data.URL.IsNull() {
		censusBaseURL = data.URL.ValueString()
	} else if censusBaseURLFromEnv, found := os.LookupEnv("CENSUS_BASE_URL"); found {
		censusBaseURL = censusBaseURLFromEnv
	}

	if censusBaseURL == "" {
		censusBaseURL = p.baseURL
	}

	if censusBaseURL == "" {
		censusBaseURL = cm.DefaultServerURL
	}

	if censusBaseURL == "" {
		resp.Diagnostics.AddAttributeError(path.Root("url"), "Failed to configure client", "No API URL provided")
	}

	var apiKey string
	if !data.ApiKey.IsNull() {
		apiKey = data.ApiKey.ValueString()
	} else {
		if apiKeyFromEnv, found := os.LookupEnv("CENSUS_WORKSPACE_API_KEY"); found {
			apiKey = apiKeyFromEnv
		}
	}

	if apiKey == "" {
		apiKey = p.apiKey
	}

	if apiKey == "" {
		resp.Diagnostics.AddAttributeError(path.Root("api_key"), "Failed to configure client", "No access token provided")
	}

	if resp.Diagnostics.HasError() {
		return
	}

	retryableClient := retryablehttp.NewClient()
	retryableClient.RetryWaitMin = time.Duration(1) * time.Second
	retryableClient.RetryWaitMax = time.Duration(3) * time.Second //nolint:mnd
	retryableClient.Backoff = retryablehttp.LinearJitterBackoff

	if p.httpClient != nil {
		retryableClient.HTTPClient = p.httpClient
	}

	censusManagementClient, err := cm.NewClient(
		censusBaseURL,
		cm.NewWorkspaceApiKeySecuritySource(apiKey),
		cm.WithClient(util.NewClientWithUserAgent(retryableClient.StandardClient(), "terraform-provider-censusworkspace/"+p.version)),
	)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create Census client: %s", err.Error())
	}

	providerData := CensusProviderData{
		client: censusManagementClient,
	}

	resp.DataSourceData = providerData
	resp.ResourceData = providerData
}

func (p *CensusProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "censusworkspace"
	resp.Version = p.version
}

func (p *CensusProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (p *CensusProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewSourceResource,
	}
}
