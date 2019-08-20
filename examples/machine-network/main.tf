provider "vra" {
    url = "${var.url}"
    access_token = "${var.access_token}"
}

resource "vra_machine" "my_machine_mysql" {
    name = "terraform_vra_mysql"
    image = "ubuntu"
    flavor = "small"

    nics {
        network_id = "${var.network_id}"
    }
}
