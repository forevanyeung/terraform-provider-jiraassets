package provider

import (
	"context"

	"github.com/ctreminiom/go-atlassian/assets"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &objectSchemaDataSource{}
	_ datasource.DataSourceWithConfigure = &objectSchemaDataSource{}
)

func NewObjectSchemaDataSource() datasource.DataSource {
	return &objectSchemaDataSource{}
}

// objectSchemaDataSource is the data source implementation.
type objectSchemaDataSource struct {
	client       *assets.Client
	workspace_id string
}

// Metadata returns the data source type name.
func (d *objectSchemaDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_object_schema"
}

// objectSchemaDataSourceModel describes the data source model.
type objectSchemaDataSourceModel struct {
	WorkspaceId     types.String `tfsdk:"workspace_id"`
	GlobalId        types.String `tfsdk:"global_id"`
	Id              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	ObjectSchemaKey types.String `tfsdk:"object_schema_key"`
	Status          types.String `tfsdk:"status"`
	Description     types.String `tfsdk:"description"`
	Created         types.String `tfsdk:"created"`
	Updated         types.String `tfsdk:"updated"`
	ObjectCount     types.Int64  `tfsdk:"object_count"`
	ObjectTypeCount types.Int64  `tfsdk:"object_type_count"`
	CanManage       types.Bool   `tfsdk:"can_manage"`
	IdAsInt         types.Int64  `tfsdk:"id_as_int"`
}

// Schema defines the schema for the data source.
func (d *objectSchemaDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"workspace_id": schema.StringAttribute{
				Computed: true,
			},
			"global_id": schema.StringAttribute{
				Computed: true,
			},
			"id": schema.StringAttribute{
				Required: true,
			},
			"name": schema.StringAttribute{
				Computed: true,
			},
			"object_schema_key": schema.StringAttribute{
				Computed: true,
			},
			"status": schema.StringAttribute{
				Computed: true,
			},
			"description": schema.StringAttribute{
				Computed: true,
			},
			"created": schema.StringAttribute{
				Computed: true,
			},
			"updated": schema.StringAttribute{
				Computed: true,
			},
			"object_count": schema.Int64Attribute{
				Computed: true,
			},
			"object_type_count": schema.Int64Attribute{
				Computed: true,
			},
			"can_manage": schema.BoolAttribute{
				Computed: true,
			},
			"id_as_int": schema.Int64Attribute{
				Computed: true,
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *objectSchemaDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Reading object schema data source")

	// Create object to hold current state of resource
	var state objectSchemaDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Call the API to get the object schema
	schema, schemaResp, err := d.client.ObjectSchema.Get(ctx, d.workspace_id, state.Id.ValueString())

	// Return an error if the API call fails
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read Assets object schema",
			err.Error(),
		)
		return
	}

	// Return an error if the HTTP status code is not 200
	if schemaResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP status code from Assets API",
			schemaResp.Status,
		)
		return
	}

	state = objectSchemaDataSourceModel{
		WorkspaceId:     types.StringValue(schema.WorkspaceId),
		GlobalId:        types.StringValue(schema.GlobalId),
		Id:              types.StringValue(schema.Id),
		Name:            types.StringValue(schema.Name),
		ObjectSchemaKey: types.StringValue(schema.ObjectSchemaKey),
		Status:          types.StringValue(schema.Status),
		Description:     types.StringValue(schema.Description),
		Created:         types.StringValue(schema.Created),
		Updated:         types.StringValue(schema.Updated),
		ObjectCount:     types.Int64Value(int64(schema.ObjectCount)),
		ObjectTypeCount: types.Int64Value(int64(schema.ObjectTypeCount)),
		CanManage:       types.BoolValue(schema.CanManage),
	}

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *objectSchemaDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerClient := req.ProviderData.(JiraAssetsProviderClient)
	// client, ok := providerClient.client.(*assets.Client)
	// if !ok {
	// 	resp.Diagnostics.AddError(
	// 		"Unexpected Data Source Configure Type",
	// 		fmt.Sprintf("Expected *assets.Client, got %T", req.ProviderData),
	// 	)
	// 	return
	// }

	d.client = providerClient.client
	d.workspace_id = providerClient.workspaceId
}
