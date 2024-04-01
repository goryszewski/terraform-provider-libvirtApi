package provider

import (
	"context"
	"fmt"

	libvirtApiClient "github.com/goryszewski/libvirtApi-client/libvirtApiClient"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type loadbalancerDataSource struct {
	client *libvirtApiClient.Client
}
type Node struct {
	name     string
	internal string
	external string
}
type Port struct {
	name     string
	protocol string
	port     types.Int64
	nodeport types.Int64
}

type loadbalancerDataSourceModel struct {
	name      string
	namespace string
	ports     []Port
	nodes     []Node
}

var (
	_ datasource.DataSource              = &loadbalancerDataSource{}
	_ datasource.DataSourceWithConfigure = &loadbalancerDataSource{}
)

func (d *loadbalancerDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = client
}
func (r *loadbalancerDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_network"
}
func (d *loadbalancerDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
func NewLoadbalancerDataSource() datasource.DataSource {
	return &loadbalancerDataSource{}
}

func (d *loadbalancerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data loadbalancerDataSourceModel

	paylaod := libvirtApiClient.LoadBalancer{
		Name:      data.name,
		Namespace: data.namespace,
	}

	diags := req.Config.Get(ctx, &data)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	loadbalancer, exist, err := d.client.GetLoadBalancer(paylaod)
	if err != nil || !exist {
		resp.Diagnostics.AddError(
			"Unable to Read loadbalancer",
			err.Error(),
		)
		return
	}
	var ports []Port
	for _, port := range loadbalancer.Ports {
		port_object := Port{
			name:     port.Name,
			protocol: port.Protocol,
			port:     types.Int64Value(int64(port.Port)),
			nodeport: types.Int64Value(int64(port.NodePort)),
		}
		ports = append(ports, port_object)
	}
	var nodes []Node
	for _, node := range loadbalancer.Nodes {
		node_object := Node{
			name:     node.Name,
			internal: node.Internal,
			external: node.External,
		}
		nodes = append(nodes, node_object)
	}

	data.nodes = nodes
	data.ports = ports

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
