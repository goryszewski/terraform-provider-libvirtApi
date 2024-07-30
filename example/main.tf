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

data "libvirtapi_network" "static" {
  id = 2
  name = "ha"
  status = 0
}

resource "libvirtapi_network" "internal01" {
  name = "ha"
}

resource "libvirtapi_network" "internal11" {
  name = "db"
}

# resource "libvirtapi_vm" "test" {
#   name = "db"
#   network = [data.libvirtapi_network.static.id,resource.libvirtapi_network.internal01.id]
# }
