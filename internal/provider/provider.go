package provider

import (
	"context"
	"net/http"
	"os"
	"time"

	cm "github.com/cysp/terraform-provider-censusworkspace/internal/census-management-go"
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

func New(version string, options ...Option) *Provider {
	provider := Provider{
		version: version,
	}

	for _, option := range options {
		option(&provider)
	}

	return &provider
}

type Provider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string

	baseURL         string
	httpClient      *http.Client
	workspaceAPIKey string
}

//nolint:revive
type ProviderModel struct {
	BaseURL         types.String `tfsdk:"base_url"`
	WorkspaceAPIKey types.String `tfsdk:"workspace_api_key"`
}

var _ provider.Provider = (*Provider)(nil)

type Option func(*Provider)

func WithBaseURL(url string) Option {
	return func(p *Provider) {
		p.baseURL = url
	}
}

func WithHTTPClient(httpClient *http.Client) Option {
	return func(p *Provider) {
		p.httpClient = httpClient
	}
}

func WithWorkspaceAPIKey(apiKey string) Option {
	return func(p *Provider) {
		p.workspaceAPIKey = apiKey
	}
}

func (p *Provider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage a Census workspace.",
		Attributes: map[string]schema.Attribute{
			"base_url": schema.StringAttribute{
				Description: "The base URL associated with your Census organization's region. If not provided, it will default to the value of the CENSUS_BASE_URL environment variable or the provider's default base URL.",
				Optional:    true,
			},
			"workspace_api_key": schema.StringAttribute{
				Description: "The API key for your Census workspace. If not provided, it will default to the value of the CENSUS_WORKSPACE_API_KEY environment variable.",
				Optional:    true,
				Sensitive:   true,
			},
		},
	}
}

func (p *Provider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data ProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var baseURL string
	if !data.BaseURL.IsNull() {
		baseURL = data.BaseURL.ValueString()
	} else if baseURLFromEnv, found := os.LookupEnv("CENSUS_BASE_URL"); found {
		baseURL = baseURLFromEnv
	}

	if baseURL == "" {
		baseURL = p.baseURL
	}

	if baseURL == "" {
		baseURL = cm.DefaultBaseURL
	}

	var workspaceAPIKey string
	if !data.WorkspaceAPIKey.IsNull() {
		workspaceAPIKey = data.WorkspaceAPIKey.ValueString()
	} else {
		if workspaceAPIKeyFromEnv, found := os.LookupEnv("CENSUS_WORKSPACE_API_KEY"); found {
			workspaceAPIKey = workspaceAPIKeyFromEnv
		}
	}

	if workspaceAPIKey == "" {
		workspaceAPIKey = p.workspaceAPIKey
	}

	if workspaceAPIKey == "" {
		resp.Diagnostics.AddAttributeError(path.Root("workspace_api_key"), "Failed to configure client", "No API Key provided")
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
		baseURL,
		cm.NewWorkspaceAPIKeySecuritySource(workspaceAPIKey),
		cm.WithClient(NewHTTPClientWithUserAgent(retryableClient.StandardClient(), "terraform-provider-censusworkspace/"+p.version)),
	)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create Census client: %s", err.Error())
	}

	providerData := ProviderData{
		client: censusManagementClient,
	}

	resp.DataSourceData = providerData
	resp.ResourceData = providerData
}

func (p *Provider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "censusworkspace"
	resp.Version = p.version
}

func (p *Provider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (p *Provider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewBigQueryDestinationResource,
		NewBigQuerySourceResource,
		NewCustomAPIDestinationResource,
		NewDestinationResource,
		NewSourceResource,
		NewSQLDatasetResource,
	}
}
