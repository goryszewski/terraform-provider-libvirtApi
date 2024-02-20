package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ provider.Provider = &libvirtapiProvider{}
)

type libvirtapiProvider struct {
	version string
}

type libvirtapiProviderModel struct {
	Hostname types.String `tfsdk:"hostname"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &libvirtapiProvider{
			version: version,
		}
	}
}

func (p *libvirtapiProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "libvirtapi"
	resp.Version = p.version
}

func (p *libvirtapiProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
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

func (p *libvirtapiProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring Libvirtapi client")
	var config libvirtapiProviderModel
	diags := req.Config.Get(ctx, &config)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.Hostname.IsUnknown() {
		resp.Diagnostics.AddAttributeError(path.Root("hostname"), "Unknown libvirtapi hostname", "...")
	}
	if config.Username.IsUnknown() {
		resp.Diagnostics.AddAttributeError(path.Root("Username"), "Unknown libvirtapi username", "...")
	}
	if config.Password.IsUnknown() {
		resp.Diagnostics.AddAttributeError(path.Root("password"), "Unknown libvirtapi password", "...")
	}

	if resp.Diagnostics.HasError() {
		return
	}

	hostname := os.Getenv("LIBVIRTapi_HOST")
	username := os.Getenv("LIBVIRTapi_USERNAME")
	password := os.Getenv("LIBVIRTapi_PASSWORD")

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
		resp.Diagnostics.AddAttributeError(path.Root("username"), "Missing libvirtapi API Username", "..")
	}
	if password == "" {
		resp.Diagnostics.AddAttributeError(path.Root("password"), "Missing libvirtapi API Password", "...")
	}

	if resp.Diagnostics.HasError() {
		return
	}
	ctx = tflog.SetField(ctx, "libvirtapi_hostname", hostname)
	ctx = tflog.SetField(ctx, "libvirtapi_username", username)
	ctx = tflog.SetField(ctx, "libvirtapi_passowrd", password)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "libvirtapi_passowrd")

	tflog.Debug(ctx, "Creating libvirtapi Client")
	// &hostname, &username, &password
	client, err := NewClient(&hostname)

	if err != nil {
		resp.Diagnostics.AddError("Unable to Create libvirtapi Client", "...")
		return
	}

	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Configured Libvirtapi client")
}

func (p *libvirtapiProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return nil
}

func (p *libvirtapiProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewNetworkResource,
	}
}
