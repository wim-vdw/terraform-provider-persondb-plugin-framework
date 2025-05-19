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
	persondbclient "github.com/wim-vdw/terraform-provider-persondb-plugin-framework/internal/client"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &persondbProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &persondbProvider{
			version: version,
		}
	}
}

// persondbProviderModel maps provider schema data to a Go type.
type persondbProviderModel struct {
	Database types.String `tfsdk:"database_filename"`
}

// persondbProvider is the provider implementation.
type persondbProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// Metadata returns the provider type name.
func (p *persondbProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "persondb"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *persondbProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"database_filename": schema.StringAttribute{
				Optional: true,
			},
		},
	}
}

// Configure prepares a PersonDB API client for data sources and resources.
func (p *persondbProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Retrieve provider data from configuration
	var config persondbProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.
	if config.Database.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("database_filename"),
			"Unknown Persons Database filename",
			"The provider cannot create the Persons DB API client as there is an unknown configuration value for the Database filename. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the CUSTOM_DATABASE_FILENAME environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.
	database := os.Getenv("CUSTOM_DATABASE_FILENAME")

	if !config.Database.IsNull() {
		database = config.Database.ValueString()
	}

	if database == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("database_filename"),
			"Missing Persons Database filename",
			"The provider cannot create the Persons DB API client as there is a missing or empty value for the Database filename. "+
				"Set the database_filename value in the configuration or use the CUSTOM_DATABASE_FILENAME environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	client, err := persondbclient.NewClient(database)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Persons DB API Client",
			"An unexpected error occurred when creating the Persons DB API Client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Persons DB API Client Error: "+err.Error(),
		)
		return
	}

	// Make the Persons DB API client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
}

// DataSources defines the data sources implemented in the provider.
func (p *persondbProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewNamesDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *persondbProvider) Resources(_ context.Context) []func() resource.Resource {
	return nil
}
