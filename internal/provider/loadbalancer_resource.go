package provider

import (
	"context"
	"fmt"

	libvirtApiClient "github.com/goryszewski/libvirtApi-client/libvirtApiClient"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &loadbalancerResource{}
	_ resource.ResourceWithConfigure   = &loadbalancerResource{}
	_ resource.ResourceWithImportState = &loadbalancerResource{}
)

// NewloadbalancerResource is a helper function to simplify the provider implementation.
func NewLoadbalancerResource() resource.Resource {
	return &loadbalancerResource{}
}

// loadbalancerResource is the resource implementation.
type loadbalancerResource struct {
	client *libvirtApiClient.Client
}

type loadbalancerResourceModel struct {
}

func (r *loadbalancerResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*libvirtApiClient.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *libvirtApiClient.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

// Metadata returns the resource type name.
func (r *loadbalancerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_loadbalancer"
}

// Schema defines the schema for the resource.
func (r *loadbalancerResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required: true,
			},
			"namespace": schema.StringAttribute{
				Required: true,
			},
			"ports": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Computed: true,
						},
						"protocol": schema.StringAttribute{
							Computed: true,
						},
						"port": schema.Int64Attribute{
							Computed: true,
						},
						"nodeport": schema.Int64Attribute{
							Computed: true,
						},
					},
				},
			},
			"nodes": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Computed: true,
						},
						"internal": schema.StringAttribute{
							Computed: true,
						},
						"external": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}

}

func (r *loadbalancerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan loadbalancerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	var bind_payload libvirtApiClient.LoadBalancer = libvirtApiClient.LoadBalancer{}

	loadbalancer, err := r.client.CreateLoadBalancer(bind_payload)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating loadbalancer",
			"Could not create loadbalancer, unexpected error: "+err.Error(),
		)
		return
	}

	plan.Name = loadbalancer

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *loadbalancerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	return
}

func (r *loadbalancerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	return
}

func (r *loadbalancerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	return
}

func (r *loadbalancerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
