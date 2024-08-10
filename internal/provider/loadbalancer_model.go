package provider

import (
	libvirtApiClient "github.com/goryszewski/libvirtApi-client/libvirtApiClient"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type loadbalancerResource struct {
	client *libvirtApiClient.Client
}

type loadbalancerDataSource struct {
	client *libvirtApiClient.Client
}

type Port struct {
	Name     string `tfsdk:"name"`
	Protocol string `tfsdk:"protocol"`
	Port     int    `tfsdk:"port"`
	NodePort int    `tfsdk:"nodeport"`
}

type Node struct {
	Name string `tfsdk:"name"`
	IP   string `tfsdk:"ip"`
}

type loadbalancerDataSourceModel struct {
	name      string
	namespace string
	ports     []Port
	nodes     []Node
}

type loadbalancerResourceModel struct {
	Ports     []Port                `tfsdk:"ports"`
	Nodes     []Node                `tfsdk:"nodes"`
	Namespace string                `tfsdk:"namespace"`
	Name      string                `tfsdk:"name"`
	Ip        basetypes.StringValue `tfsdk:"ip"`
}
