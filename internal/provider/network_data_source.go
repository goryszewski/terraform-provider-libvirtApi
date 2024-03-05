package provider

import (
	"context"
	"fmt"
	"strconv"

	libvirtApiClient "github.com/goryszewski/libvirtApi-client/libvirtApiClient"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type networkDataSource struct {
	client *libvirtApiClient.Client
}
type networkDataSourceModel struct {
	ID     types.Int64 `tfsdk:"id"`
	Name   string      `tfsdk:"name"`
	Status types.Int64 `tfsdk:"status"`
}

var (
	_ datasource.DataSource              = &networkDataSource{}
	_ datasource.DataSourceWithConfigure = &networkDataSource{}
)

// Configure adds the provider configured client to the data source.
func (d *networkDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// Metadata returns the resource type name.
func (r *networkDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_network"
}

func NewNetworkDataSource() datasource.DataSource {
	return &networkDataSource{}
}

func (d *networkDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Required: true,
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"status": schema.Int64Attribute{
				Required: true,
			},
		},
	}
}
func (d *networkDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data networkDataSourceModel

	diags := req.Config.Get(ctx, &data)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	network, err := d.client.GetNetwork(int(data.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Network ",
			err.Error(),
		)
		return
	}

	if network.Name != data.Name {

		data_string := network.Name + " != " + data.Name
		resp.Diagnostics.AddError(
			"Bad name:", data_string,
		)
		return
	}

	if types.Int64Value(int64(network.Status)) != data.Status {

		data_string := strconv.Itoa(network.Status) + " != " + data.Status.String()
		resp.Diagnostics.AddError(
			"Bad Status:", data_string,
		)
		return
	}

	// data.ID = types.Int64Value(int64(network.ID))
	// data.Name = network.Name
	// data.Status = types.Int64Value(int64(network.Status))

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
