package provider

import (
	"context"
	"fmt"

	libvirtApiClient "github.com/goryszewski/libvirtApi-client/libvirtApiClient"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
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
			"ip": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Required: true,

				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"namespace": schema.StringAttribute{
				Required: true,

				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"ports": schema.ListNestedAttribute{
				Required: true,

				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required: true,

							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
						},
						"protocol": schema.StringAttribute{
							Required: true,
						},
						"port": schema.Int64Attribute{
							Required: true,
						},
						"nodeport": schema.Int64Attribute{
							Required: true,
						},
					},
				},
			},
			"nodes": schema.ListNestedAttribute{
				Required: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required: true,
						},
						"ip": schema.StringAttribute{
							Required: true,
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

	var bind_payload libvirtApiClient.LoadBalancer = libvirtApiClient.LoadBalancer{
		Name:      plan.Name,
		Namespace: plan.Namespace,
	}
	for _, node := range plan.Nodes {
		var tmp libvirtApiClient.Node = libvirtApiClient.Node{
			Name: node.Name,
			IP:   node.IP,
		}
		bind_payload.Nodes = append(bind_payload.Nodes, tmp)
	}
	for _, port := range plan.Ports {
		var tmp libvirtApiClient.Port_Service = libvirtApiClient.Port_Service{
			Name:     port.Name,
			Protocol: port.Protocol,
			Port:     port.Port,
			NodePort: port.NodePort,
		}
		bind_payload.Ports = append(bind_payload.Ports, tmp)
	}

	ip, err := r.client.CreateLoadBalancer(bind_payload)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating loadbalancer",
			"Could not create loadbalancer, unexpected error: "+err.Error(),
		)
		return
	}

	plan.Ip = basetypes.NewStringValue(ip)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *loadbalancerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state loadbalancerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var bind_payload libvirtApiClient.LoadBalancer = libvirtApiClient.LoadBalancer{
		Name:      state.Name,
		Namespace: state.Namespace,
	}
	for _, node := range state.Nodes {
		var tmp libvirtApiClient.Node = libvirtApiClient.Node{
			Name: node.Name,
			IP:   node.IP,
		}
		bind_payload.Nodes = append(bind_payload.Nodes, tmp)
	}
	for _, port := range state.Ports {
		var tmp libvirtApiClient.Port_Service = libvirtApiClient.Port_Service{
			Name:     port.Name,
			Protocol: port.Protocol,
			Port:     port.Port,
			NodePort: port.NodePort,
		}
		bind_payload.Ports = append(bind_payload.Ports, tmp)
	}

	lb, _, err := r.client.GetLoadBalancer(bind_payload)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading lb",
			"Could not read lb ID : "+err.Error(),
		)
		return
	}

	state.Ip = basetypes.NewStringValue(lb.Ip)
	state.Name = lb.Name
	state.Namespace = lb.Namespace
	state.Nodes = []Node{}
	for _, node := range lb.Nodes {

		state.Nodes = append(state.Nodes, Node{
			Name: node.Name,
			IP:   node.IP,
		})
	}
	state.Ports = []Port{}
	for _, node := range lb.Ports {

		state.Ports = append(state.Ports, Port{
			Name:     node.Name,
			Protocol: node.Protocol,
			Port:     node.Port,
			NodePort: node.NodePort,
		})
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *loadbalancerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state loadbalancerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	var plan loadbalancerResourceModel
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	var bind_payload libvirtApiClient.LoadBalancer = libvirtApiClient.LoadBalancer{
		Name:      plan.Name,
		Namespace: plan.Namespace,
	}
	for _, node := range plan.Nodes {
		var tmp libvirtApiClient.Node = libvirtApiClient.Node{
			Name: node.Name,
			IP:   node.IP,
		}
		bind_payload.Nodes = append(bind_payload.Nodes, tmp)
	}
	for _, port := range plan.Ports {
		var tmp libvirtApiClient.Port_Service = libvirtApiClient.Port_Service{
			Name:     port.Name,
			Protocol: port.Protocol,
			Port:     port.Port,
			NodePort: port.NodePort,
		}
		bind_payload.Ports = append(bind_payload.Ports, tmp)
	}

	err := r.client.UpdateLoadBalancer(bind_payload)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error update loadbalancer",
			"Could not update loadbalancer, unexpected error: "+err.Error(),
		)
		return
	}
	state.Name = bind_payload.Name
	state.Namespace = bind_payload.Namespace

	state.Nodes = []Node{}
	for _, node := range bind_payload.Nodes {

		state.Nodes = append(state.Nodes, Node{
			Name: node.Name,
			IP:   node.IP,
		})
	}
	state.Ports = []Port{}
	for _, node := range bind_payload.Ports {

		state.Ports = append(state.Ports, Port{
			Name:     node.Name,
			Protocol: node.Protocol,
			Port:     node.Port,
			NodePort: node.NodePort,
		})
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *loadbalancerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state loadbalancerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	var bind_payload libvirtApiClient.LoadBalancer = libvirtApiClient.LoadBalancer{
		Name:      state.Name,
		Namespace: state.Namespace,
	}
	for _, node := range state.Nodes {
		var tmp libvirtApiClient.Node = libvirtApiClient.Node{
			Name: node.Name,
			IP:   node.IP,
		}
		bind_payload.Nodes = append(bind_payload.Nodes, tmp)
	}
	for _, port := range state.Ports {
		var tmp libvirtApiClient.Port_Service = libvirtApiClient.Port_Service{
			Name:     port.Name,
			Protocol: port.Protocol,
			Port:     port.Port,
			NodePort: port.NodePort,
		}
		bind_payload.Ports = append(bind_payload.Ports, tmp)
	}
	err := r.client.DeleteLoadBalancer(bind_payload)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting LoadBalancer",
			"Could not delete LoadBalancer, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *loadbalancerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
