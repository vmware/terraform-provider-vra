provider "cas" {
    url = "${var.url}"
    access_token = "${var.access_token}"
    project_id = "${var.project_id}"
    deployment_id = "${var.deployment_id}"
}

resource "cas_machine" "my_machine_mysql" {
    name = "terraform_cas_mysql"
    image = "ubuntu"
    flavor = "small"

    nics {
        network_id = "${var.network_id}"
    }
}