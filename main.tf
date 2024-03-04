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

resource "libvirtapi_network" "internal01" {
  name = "ha"
}

resource "libvirtapi_network" "internal11" {
  name = "db"
}
