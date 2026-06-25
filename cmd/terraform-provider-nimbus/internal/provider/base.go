package provider

import (
	"context"
	"fmt"

	"github.com/leetrout/terraform-sim/cmd/terraform-provider-nimbus/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// baseResource implements the shared CRUD/import behavior for every nimbus
// resource. The per-resource files declare a spec describing the Terraform
// type name, REST collection, and writable string attributes; this base does
// the rest, mapping generic JSON objects to/from Terraform state.
//
// State is modeled as a dynamic object: a fixed set of attributes (the
// computed "id" and "serial" plus the writable strings declared in the spec).
type baseResource struct {
	spec   resourceSpec
	client *client.Client
}

// resourceSpec describes one nimbus resource type.
type resourceSpec struct {
	// typeName is the full Terraform resource type, e.g. "nimbus_network".
	typeName string
	// collection is the REST collection path segment, e.g. "networks".
	collection string
	// description is shown in docs.
	description string
	// writable lists the writable string attributes in declaration order.
	writable []attrSpec
}

type attrSpec struct {
	name        string
	description string
	// requiresReplace forces resource replacement when changed (e.g. fields the
	// API treats as immutable parent references).
	requiresReplace bool
	optional        bool
}

// objectAttrTypes returns the attr.Type map for this resource's state object.
func (s resourceSpec) objectAttrTypes() map[string]attr.Type {
	m := map[string]attr.Type{
		"id":     types.StringType,
		"serial": types.Int64Type,
	}
	for _, a := range s.writable {
		m[a.name] = types.StringType
	}
	return m
}

// Schema builds the Terraform schema from the spec.
func (r *baseResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	attrs := map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "Server-assigned identifier.",
			Computed:            true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"serial": schema.Int64Attribute{
			MarkdownDescription: "Monotonic revision counter, bumped on every mutation.",
			Computed:            true,
		},
	}
	for _, a := range r.spec.writable {
		var mods []planmodifier.String
		if a.requiresReplace {
			mods = append(mods, stringplanmodifier.RequiresReplace())
		}
		attrs[a.name] = schema.StringAttribute{
			MarkdownDescription: a.description,
			Required:            !a.optional,
			Optional:            a.optional,
			PlanModifiers:       mods,
		}
	}

	resp.Schema = schema.Schema{
		MarkdownDescription: r.spec.description,
		Attributes:          attrs,
	}
}

func (r *baseResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = r.spec.typeName
}

func (r *baseResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Provider Data Type",
			fmt.Sprintf("Expected *client.Client, got %T. This is a provider bug.", req.ProviderData),
		)
		return
	}
	r.client = c
}

// payloadFromState extracts the writable attributes from a state object into a
// generic JSON map for sending to the API.
func (r *baseResource) payloadFromState(ctx context.Context, obj types.Object) (map[string]any, error) {
	attrs := obj.Attributes()
	payload := make(map[string]any, len(r.spec.writable))
	for _, a := range r.spec.writable {
		v, ok := attrs[a.name].(types.String)
		if !ok {
			return nil, fmt.Errorf("attribute %q missing or not a string", a.name)
		}
		if v.IsNull() {
			continue
		}
		payload[a.name] = v.ValueString()
	}
	return payload, nil
}

// stateFromResponse builds a Terraform object value from an API response map.
func (r *baseResource) stateFromResponse(resp map[string]any) (types.Object, error) {
	values := map[string]attr.Value{}

	id, _ := resp["id"].(string)
	values["id"] = types.StringValue(id)

	// JSON numbers decode to float64.
	var serial int64
	switch v := resp["serial"].(type) {
	case float64:
		serial = int64(v)
	case int64:
		serial = v
	}
	values["serial"] = types.Int64Value(serial)

	for _, a := range r.spec.writable {
		s, _ := resp[a.name].(string)
		values[a.name] = types.StringValue(s)
	}

	obj, diags := types.ObjectValue(r.spec.objectAttrTypes(), values)
	if diags.HasError() {
		return types.ObjectNull(r.spec.objectAttrTypes()), fmt.Errorf("build state object: %v", diags)
	}
	return obj, nil
}

func (r *baseResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan types.Object
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	payload, err := r.payloadFromState(ctx, plan)
	if err != nil {
		resp.Diagnostics.AddError("Invalid Plan", err.Error())
		return
	}

	var out map[string]any
	if err := r.client.Create(r.spec.collection, payload, &out); err != nil {
		resp.Diagnostics.AddError("Create Failed", err.Error())
		return
	}

	state, err := r.stateFromResponse(out)
	if err != nil {
		resp.Diagnostics.AddError("Invalid Response", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *baseResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state types.Object
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, _ := state.Attributes()["id"].(types.String)
	if id.IsNull() || id.ValueString() == "" {
		resp.State.RemoveResource(ctx)
		return
	}

	var out map[string]any
	if err := r.client.Get(r.spec.collection, id.ValueString(), &out); err != nil {
		if client.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Read Failed", err.Error())
		return
	}

	newState, err := r.stateFromResponse(out)
	if err != nil {
		resp.Diagnostics.AddError("Invalid Response", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *baseResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan types.Object
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	var state types.Object
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, _ := state.Attributes()["id"].(types.String)

	payload, err := r.payloadFromState(ctx, plan)
	if err != nil {
		resp.Diagnostics.AddError("Invalid Plan", err.Error())
		return
	}

	var out map[string]any
	if err := r.client.Update(r.spec.collection, id.ValueString(), payload, &out); err != nil {
		resp.Diagnostics.AddError("Update Failed", err.Error())
		return
	}

	newState, err := r.stateFromResponse(out)
	if err != nil {
		resp.Diagnostics.AddError("Invalid Response", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *baseResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state types.Object
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, _ := state.Attributes()["id"].(types.String)
	if id.IsNull() || id.ValueString() == "" {
		return
	}

	if err := r.client.Delete(r.spec.collection, id.ValueString()); err != nil {
		if client.IsNotFound(err) {
			return
		}
		resp.Diagnostics.AddError("Delete Failed", err.Error())
		return
	}
}

func (r *baseResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Compile-time assertions that baseResource satisfies the framework interfaces.
var (
	_ resource.Resource                = (*baseResource)(nil)
	_ resource.ResourceWithConfigure   = (*baseResource)(nil)
	_ resource.ResourceWithImportState = (*baseResource)(nil)
)
