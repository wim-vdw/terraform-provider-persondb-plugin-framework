package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	persondbclient "github.com/wim-vdw/terraform-provider-persondb-plugin-framework/internal/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &PersonResource{}
var _ resource.ResourceWithImportState = &PersonResource{}

// NewPersonResource is a helper function to simplify the provider implementation.
func NewPersonResource() resource.Resource {
	return &PersonResource{}
}

// PersonResource is the resource implementation.
type PersonResource struct {
	client *persondbclient.Client
}

// PersonResourceModel maps the resource schema data.
type PersonResourceModel struct {
	ID        types.String `tfsdk:"id"`
	PersonID  types.String `tfsdk:"person_id"`
	LastName  types.String `tfsdk:"last_name"`
	FirstName types.String `tfsdk:"first_name"`
}

// Metadata returns the resource type name.
func (r *PersonResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_person"
}

// Schema defines the schema for the resource.
func (r *PersonResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"person_id": schema.StringAttribute{
				Required: true,
			},
			"first_name": schema.StringAttribute{
				Optional: true,
			},
			"last_name": schema.StringAttribute{
				Required: true,
			},
		},
	}
}

func (r *PersonResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*persondbclient.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *persondbclient.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *PersonResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data PersonResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	personID := data.PersonID.ValueString()
	lastName := data.LastName.ValueString()
	firstName := data.FirstName.ValueString()

	// Check if the person already exists
	exists, err := r.client.CheckPersonExists(personID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating person",
			"Could not create person, unexpected error: "+err.Error(),
		)
		return
	}
	if exists {
		resp.Diagnostics.AddError(
			"Error creating person",
			"Person with person_id '"+personID+"' already exists. Use 'terraform import' to manage it in Terraform.",
		)
		return
	}

	// Create new person
	err = r.client.CreatePerson(personID, lastName, firstName)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating person",
			"Could not create person, unexpected error: "+err.Error(),
		)
		return
	}

	// Save ID with the format "/person/<person_id>" to Terraform state
	data.ID = types.StringValue("/person/" + personID)

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PersonResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data PersonResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	parts := strings.Split(data.ID.ValueString(), "/")
	if len(parts) != 3 || parts[1] != "person" {
		resp.Diagnostics.AddError(
			"Invalid ID format",
			fmt.Sprintf("Expected '/person/<person_id>', got: %s", data.ID.ValueString()),
		)
		return
	}

	personID := parts[2]
	lastName, firstName, err := r.client.ReadPerson(personID)
	if err != nil {
		// Person could not be found, so we call the RemoveResource method to enforce new resource creation
		resp.State.RemoveResource(ctx)
		return
	}

	data.PersonID = types.StringValue(personID)
	data.LastName = types.StringValue(lastName)
	data.FirstName = types.StringValue(firstName)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PersonResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data PersonResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	personID := data.PersonID.ValueString()
	lastName := data.LastName.ValueString()
	firstName := data.FirstName.ValueString()

	// Update person
	err := r.client.UpdatePerson(personID, lastName, firstName)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating person",
			"Could not update person, unexpected error: "+err.Error(),
		)
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PersonResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data PersonResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	personID := data.PersonID.ValueString()

	// Delete person
	err := r.client.DeletePerson(personID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting person",
			"Could not delete person, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *PersonResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
