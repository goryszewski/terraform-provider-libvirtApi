package provider

import (
	"context"
	"os"

	libvirtApiClient "github.com/goryszewski/libvirtApi-client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ provider.Provider = &libvirtApiProvider{}
)

type libvirtApiProvider struct {
	version string
}

type libvirtApiProviderModel struct {
	Hostname types.String `trsddk:"hostname"`
	Username types.String `trsddk:"username"`
	Password types.String `trsddk:"password"`
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &libvirtApiProvider{
			version: version,
		}
	}
}

func (p *libvirtApiProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "libvirtApi"
	resp.Version = p.version
}

func (p *libvirtApiProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"hostname": schema.StringAttribute{
				Optional: true,
			},
			"username": schema.StringAttribute{
				Optional: true,
			},
			"password": schema.StringAttribute{
				Optional: true,
			},
		},
	}
}

func (p *libvirtApiProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring LibvirtApi client")
	var config libvirtApiProviderModel
	diags := req.Config.Get(ctx, &config)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.Hostname.IsUnknown() {
		resp.Diagnostics.AddAttributeError(path.Root("hostname"), "Unknown libvirtApi hostname", "...")
	}
	if config.Username.IsUnknown() {
		resp.Diagnostics.AddAttributeError(path.Root("Username"), "Unknown libvirtApi username", "...")
	}
	if config.Password.IsUnknown() {
		resp.Diagnostics.AddAttributeError(path.Root("password"), "Unknown libvirtApi password", "...")
	}

	if resp.Diagnostics.HasError() {
		return
	}

	hostname := os.Getenv("LIBVIRTAPI_HOST")
	username := os.Getenv("LIBVIRTAPI_USERNAME")
	password := os.Getenv("LIBVIRTAPI_PASSWORD")

	if !config.Hostname.IsNull() {
		hostname = config.Hostname.ValueString()
	}
	if !config.Username.IsNull() {
		username = config.Hostname.ValueString()
	}
	if !config.Password.IsNull() {
		password = config.Hostname.ValueString()
	}

	if hostname == "" {
		resp.Diagnostics.AddAttributeError(path.Root("hostname"), "Missing livbirtApi Hostname", "...")
	}
	if username == "" {
		resp.Diagnostics.AddAttributeError(path.Root("username"), "Missing libvirtApi API Username", "..")
	}
	if password == "" {
		resp.Diagnostics.AddAttributeError(path.Root("password"), "Missing libvirtApi API Password", "...")
	}

	if resp.Diagnostics.HasError() {
		return
	}
	ctx = tflog.SetField(ctx, "libvirtApi_hostname", hostname)
	ctx = tflog.SetField(ctx, "libvirtApi_username", username)
	ctx = tflog.SetField(ctx, "libvirtApi_passowrd", password)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "libvirtApi_passowrd")

	tflog.Debug(ctx, "Creating libvirtApi Client")
	client, err := libvirtApiClient.NewClient(&hostname, &username, &password)

	if err != nil {
		resp.Diagnostics.AddError("Unable to Create libvirtApi Client", "...")
		return
	}

	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Configured LibvirtApi client")
}

func (p *libvirtApiProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return nil
}

func (p *libvirtApiProvider) Resources(_ context.Context) []func() resource.Resource {
	return nil
}
