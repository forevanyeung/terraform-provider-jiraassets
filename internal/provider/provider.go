// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/ctreminiom/go-atlassian/assets"
)

// Ensure the implementation satisfies the expected interfaces.
var _ provider.Provider = &JiraAssetsProvider{}

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &JiraAssetsProvider{
			version: version,
		}
	}
}

// JiraAssetsProvider defines the provider implementation.
type JiraAssetsProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// JiraAssetsProviderModel describes the provider data model.
type JiraAssetsProviderModel struct {
	WorkspaceId types.String `tfsdk:"workspace_id"`
	User        types.String `tfsdk:"user"`
	Password    types.String `tfsdk:"password"`
}

// struct describes client and worksapceId
type JiraAssetsProviderClient struct {
	client *assets.Client
	workspaceId string
}

func (p *JiraAssetsProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "jiraassets"
	resp.Version = p.version
}

func (p *JiraAssetsProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"workspace_id": schema.StringAttribute{
				MarkdownDescription: "Workspace Id of the Asssets instance.",
				Optional:            true,
			},
			"user": schema.StringAttribute{
				MarkdownDescription: "Username of an admin or service account with access to the Jira API.",
				Optional:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "Personal access token for the admin or service account.",
				Optional:            true,
				Sensitive:           true,
			},
		},
	}
}

func (p *JiraAssetsProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config JiraAssetsProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the attributes, it must be a known value.

	if config.WorkspaceId.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("workspaceId"),
			"Unknown Assets Workspace Id",
			"The provider cannot create the Assets API client as there is an unknown configuration value for the Assets API workspace Id. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the JIRASSETS_WORKSPACE_ID environment variable.",
		)
	}

	if config.User.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("user"),
			"Unknown Assets User",
			"The provider cannot create the Assets API client as there is an unknown configuration value for the Assets API user. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the JIRASSETS_USER environment variable.",
		)
	}

	if config.Password.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Unknown Assets Password",
			"The provider cannot create the Assets API client as there is an unknown configuration value for the Assets API password. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the JIRASSETS_PASSWORD environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override with Terraform configuration value if set.

	workspaceId := os.Getenv("JIRASSETS_WORKSPACE_ID")
	user := os.Getenv("JIRASSETS_USER")
	password := os.Getenv("JIRASSETS_PASSWORD")

	if !config.WorkspaceId.IsNull() {
		workspaceId = config.WorkspaceId.ValueString()
	}

	if !config.User.IsNull() {
		user = config.User.ValueString()
	}

	if !config.Password.IsNull() {
		password = config.Password.ValueString()
	}

	// If any of the expected configurations are missing, return errors with provider-specific guidance.

	if workspaceId == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("workspaceId"),
			"Missing Assets API Workspace Id",
			"The provider cannot create the Assets API client as there is a missing or empty value for the Assets API workspace Id. "+
				"Set the host value in the configuration or use the JIRASSETS_WORKSPACE_ID environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if user == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("user"),
			"Missing Assets API User",
			"The provider cannot create the Assets API client as there is a missing or empty value for the Assets API username. "+
				"Set the user value in the configuration or use the JIRASSETS_USER environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if password == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Missing Assets API Password",
			"The provider cannot create the Assets API client as there is a missing or empty value for the Assets API password. "+
				"Set the password value in the configuration or use the JIRASSETS_PASSWORD environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// create the Jira Assets client
	client, err := assets.New(nil, "")
	
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create Assets client",
			"An unexpected error occured wehn creating the Assets API client. Error: "+err.Error(),
		)
	}
	
	// add authentication headers to the client, workspaceId is added to each request
	client.Auth.SetBasicAuth(user, password)

	// add workspaceId to response to be used by resources and data sources
	providerClient := JiraAssetsProviderClient{
		client: client,
		workspaceId: workspaceId,
	}

	resp.DataSourceData = providerClient
	resp.ResourceData = providerClient
}

func (p *JiraAssetsProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewExampleResource,
	}
}

func (p *JiraAssetsProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewExampleDataSource,
		NewObjectSchemaDataSource,
	}
}
