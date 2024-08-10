terraform {
  required_providers {
    libvirtapi = {
      source = "github.com/goryszewski/libvirtApi"
    }
  }
  # required_version = ">= 1.1.0"
}

provider "libvirtapi" {
  hostname = "http://127.0.0.1:8050"
  username = "test"
  password = "test"
}

resource "libvirtapi_loadbalancer" "lbApi" {
  name = "db"
	namespace= "ee"
  nodes = [{
    name = "12"
    ip = "3.3.2.121"
  },{
    name = "13"
    ip = "3.3.2.13"
  }
  ]
  ports = [{
    name = "test"
    protocol = "tcp"
    port = "801"
    nodeport = "1234"
  },
  {
    name = "test11"
    protocol = "tcp"
    port = "81"
    nodeport = "123"
  }
  ]
}


# data "libvirtapi_network" "static" {
#   id = 2
#   name = "ha"
#   status = 0
# }

# resource "libvirtapi_network" "internal01" {
#   name = "ha"
# }

# resource "libvirtapi_network" "internal11" {
#   name = "db"
# }

# resource "libvirtapi_vm" "test" {
#   name = "db"
#   network = [data.libvirtapi_network.static.id,resource.libvirtapi_network.internal01.id]
# }
