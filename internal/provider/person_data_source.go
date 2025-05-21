package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	persondbclient "github.com/wim-vdw/terraform-provider-persondb-plugin-framework/internal/client"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource = &PersonDataSource{}
)

// NewPersonDataSource is a helper function to simplify the provider implementation.
func NewPersonDataSource() datasource.DataSource {
	return &PersonDataSource{}
}

// PersonDataSource is the data source implementation.
type PersonDataSource struct {
	client *persondbclient.Client
}

// PersonDataSourceModel maps the data source schema data.
type PersonDataSourceModel struct {
	ID        types.String `tfsdk:"id"`
	PersonID  types.String `tfsdk:"person_id"`
	LastName  types.String `tfsdk:"last_name"`
	FirstName types.String `tfsdk:"first_name"`
}

// Metadata returns the data source type name.
func (d *PersonDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_person"
}

// Schema defines the schema for the data source.
func (d *PersonDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"person_id": schema.StringAttribute{
				Required: true,
			},
			"last_name": schema.StringAttribute{
				Computed: true,
			},
			"first_name": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *PersonDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data PersonDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	personId := data.PersonID.ValueString()
	lastName, firstName, err := d.client.ReadPerson(personId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading person",
			"Could not read person, unexpected error: "+err.Error(),
		)
		return
	}

	data.ID = types.StringValue("/person/" + personId)
	data.LastName = types.StringValue(lastName)

	// Only set first name if it is not empty
	if firstName != "" {
		data.FirstName = types.StringValue(firstName)
	}

	// Set data
	diags := resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *PersonDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*persondbclient.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *persondbclient.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}
