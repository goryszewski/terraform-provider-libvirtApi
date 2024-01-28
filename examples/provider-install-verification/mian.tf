terraform {
  required_providers {
    hashicups = {
      source = "github.com/goryszewski/libvirtApi"
    }
  }
}

provider "hashicups" {
  hostname = "http://127.0.0.1:8050"
  username = "test"
  password = "test"
}


