package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource = &namesDataSource{}
)

// namesDataSourceModel maps the data source schema data.
type namesDataSourceModel struct {
	Names []namesModel `tfsdk:"names"`
}

// namesModel maps names schema data.
type namesModel struct {
	FirstName types.String `tfsdk:"first_name"`
	LastName  types.String `tfsdk:"last_name"`
}

// NewNamesDataSource is a helper function to simplify the provider implementation.
func NewNamesDataSource() datasource.DataSource {
	return &namesDataSource{}
}

// namesDataSource is the data source implementation.
type namesDataSource struct{}

// Metadata returns the data source type name.
func (d *namesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_names"
}

// Schema defines the schema for the data source.
func (d *namesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"names": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"first_name": schema.StringAttribute{
							Computed: true,
						},
						"last_name": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *namesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state namesDataSourceModel

	wim := namesModel{
		FirstName: types.StringValue("Wim"),
		LastName:  types.StringValue("Van den Wyngaert"),
	}

	state.Names = append(state.Names, wim)

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
