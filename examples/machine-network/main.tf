provider "vra" {
  url          = var.url
  access_token = var.access_token
}

resource "vra_machine" "this" {
  name   = "terraform-vra-machine"
  image  = "ubuntu"
  flavor = "small"

  nics {
    network_id = var.network_id
  }
}
