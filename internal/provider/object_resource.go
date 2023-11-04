package provider

import (
	"context"

	"github.com/ctreminiom/go-atlassian/assets"
	"github.com/ctreminiom/go-atlassian/pkg/infra/models"

	// "github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &objectResource{}
	_ resource.ResourceWithConfigure = &objectResource{}
)

// NewObjectResource is a helper function to simplify the provider implementation.
func NewObjectResource() resource.Resource {
	return &objectResource{}
}

// objectResource is the resource implementation.
type objectResource struct {
	client       *assets.Client
	workspace_id string
}

// Metadata returns the resource type name.
func (r *objectResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_object"
}

type objectResourceModel struct {

	// avatar
	// objectType
	// hasAvatar
	// timestamp
	// attributes
	// _links

	WorkspaceId types.String `tfsdk:"workspace_id"`
	GlobalId    types.String `tfsdk:"global_id"`
	Id          types.String `tfsdk:"id"`
	Label       types.String `tfsdk:"label"`
	ObjectKey   types.String `tfsdk:"object_key"`
	Created     types.String `tfsdk:"created"`
	Updated     types.String `tfsdk:"updated"`
	HasAvatar   types.Bool   `tfsdk:"has_avatar"`

	TypeId     types.String              `tfsdk:"type_id"`
	Attributes []objectAttrResourceModel `tfsdk:"attributes"`
	AvatarUuid types.String              `tfsdk:"avatar_uuid"`
}

type objectAttrResourceModel struct {
	AttrTypeId types.String `tfsdk:"attr_type_id"`
	AttrValue  types.String `tfsdk:"attr_value"`
}

// Schema defines the schema for the resource.
func (r *objectResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"workspace_id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"global_id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"label": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"object_key": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"type_id": schema.StringAttribute{
				Required: true,
			},
			"attributes": schema.SetNestedAttribute{
				Required: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"attr_type_id": schema.StringAttribute{
							Required: true,
						},
						"attr_value": schema.StringAttribute{
							Required: true,
						},
					},
				},
			},
			"created": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"updated": schema.StringAttribute{
				Computed: true,
			},
			"has_avatar": schema.BoolAttribute{
				Optional: true,
			},
			"avatar_uuid": schema.StringAttribute{
				Optional: true,
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *objectResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan objectResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var attributes []*models.ObjectPayloadAttributeScheme
	for _, attr := range plan.Attributes {
		attributes = append(attributes, &models.ObjectPayloadAttributeScheme{
			ObjectTypeAttributeID: attr.AttrTypeId.ValueString(),
			ObjectAttributeValues: []*models.ObjectPayloadAttributeValueScheme{
				{
					Value: attr.AttrValue.ValueString(),
				},
			},
		})
	}

	// create payload
	payload := &models.ObjectPayloadScheme{
		ObjectTypeID: plan.TypeId.ValueString(),
		Attributes:   attributes,
		HasAvatar:    plan.HasAvatar.ValueBool(),
		AvatarUUID:   plan.AvatarUuid.ValueString(),
	}

	object, response, err := r.client.Object.Create(ctx, r.workspace_id, payload)
	if err != nil {
		if response != nil {
			tflog.Error(ctx, "Error creating object: %s", map[string]interface{}{
				"url":         response.Request.URL,
				"status_code": response.StatusCode,
				"headers":     response.Header,
				"body":        response.Body,
			})
		}

		resp.Diagnostics.AddError(
			"Error during object creation",
			err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attributes
	plan.WorkspaceId = types.StringValue(object.WorkspaceId)
	plan.GlobalId = types.StringValue(object.GlobalId)
	plan.Id = types.StringValue(object.ID)
	plan.Label = types.StringValue(object.Label)
	plan.ObjectKey = types.StringValue(object.ObjectKey)
	plan.Created = types.StringValue(object.Created)
	plan.Updated = types.StringValue(object.Updated)
	plan.HasAvatar = types.BoolValue(object.HasAvatar)

	// Set state to full populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *objectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state objectResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed object from Assets API
	object, response, err := r.client.Object.Get(ctx, r.workspace_id, state.Id.ValueString())
	if err != nil {
		if response != nil {
			tflog.Error(ctx, "Error reading object: %s", map[string]interface{}{
				"url":         response.Request.URL,
				"status_code": response.StatusCode,
				"headers":     response.Header,
				"body":        response.Body,
			})
		}

		resp.Diagnostics.AddError(
			"Error during object reading",
			err.Error(),
		)
		return
	}

	// Get refreshed object attributes from Assets API
	attrs, response, err := r.client.Object.Attributes(ctx, r.workspace_id, state.Id.ValueString())
	if err != nil {
		if response != nil {
			tflog.Error(ctx, "Error reading object attributes: %s", map[string]interface{}{
				"url":         response.Request.URL,
				"status_code": response.StatusCode,
				"headers":     response.Header,
				"body":        response.Body,
			})
		}
		resp.Diagnostics.AddError(
			"Error during object attributes reading",
			err.Error(),
		)
		return
	}

	var attributes []objectAttrResourceModel
	for _, attr := range attrs {
		// only map known attributes in the state, this is because the API return computed attributes like "key", "created",
		// and "updated". we don't know the type id of those attributes, so we can't exclude them specifically

		for i := range state.Attributes {
			if state.Attributes[i].AttrTypeId == types.StringValue(attr.ObjectTypeAttributeId) {
				attributes = append(attributes, objectAttrResourceModel{
					AttrTypeId: types.StringValue(attr.ObjectTypeAttributeId),
					AttrValue:  types.StringValue(attr.ObjectAttributeValues[0].Value),
				})
			}
		}
	}

	// Overwrite items in state with refreshed values
	state.Attributes = attributes
	state.WorkspaceId = types.StringValue(object.WorkspaceId)
	state.GlobalId = types.StringValue(object.GlobalId)
	state.Id = types.StringValue(object.ID)
	state.Label = types.StringValue(object.Label)
	state.ObjectKey = types.StringValue(object.ObjectKey)
	state.HasAvatar = types.BoolValue(object.HasAvatar)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *objectResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan objectResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	// if an attribute is removed from plan, it will not be removed from the object
	// this is due to how the API only partially updates the object
	var attributes []*models.ObjectPayloadAttributeScheme
	for _, attr := range plan.Attributes {
		attributes = append(attributes, &models.ObjectPayloadAttributeScheme{
			ObjectTypeAttributeID: attr.AttrTypeId.ValueString(),
			ObjectAttributeValues: []*models.ObjectPayloadAttributeValueScheme{
				{
					Value: attr.AttrValue.ValueString(),
				},
			},
		})
	}

	// create payload
	payload := &models.ObjectPayloadScheme{
		ObjectTypeID: plan.TypeId.ValueString(),
		Attributes:   attributes,
		HasAvatar:    plan.HasAvatar.ValueBool(),
		AvatarUUID:   plan.AvatarUuid.ValueString(),
	}

	// update object
	tflog.Info(ctx, "Updating object.", map[string]interface{}{
		"Id": plan.Id.ValueString(),
	})
	object, response, err := r.client.Object.Update(ctx, r.workspace_id, plan.Id.ValueString(), payload)
	if err != nil {
		if response != nil {
			tflog.Error(ctx, "Error updating object: %s", map[string]interface{}{
				"url":         response.Request.URL,
				"status_code": response.StatusCode,
				"headers":     response.Header,
				"body":        response.Body,
			})
		}

		resp.Diagnostics.AddError(
			"Error during object update",
			err.Error(),
		)
		return
	}

	// Update resource state with updated object and attributes
	plan.WorkspaceId = types.StringValue(object.WorkspaceId)
	plan.GlobalId = types.StringValue(object.GlobalId)
	plan.Id = types.StringValue(object.ID)
	plan.Label = types.StringValue(object.Label)
	plan.ObjectKey = types.StringValue(object.ObjectKey)
	plan.Created = types.StringValue(object.Created)
	plan.Updated = types.StringValue(object.Updated)
	plan.HasAvatar = types.BoolValue(object.HasAvatar)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *objectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

// Configure configures the resource with the given configuration.
func (r *objectResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerClient := req.ProviderData.(JiraAssetsProviderClient)

	r.client = providerClient.client
	r.workspace_id = providerClient.workspaceId
}
