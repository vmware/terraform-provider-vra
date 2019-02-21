provider "tango" {
    url = "${var.url}"
    access_token = "${var.access_token}"
    project_id = "${var.project_id}"
    deployment_id = "${var.deployment_id}"
}

resource "tango_machine" "my_machine_mysql" {
    name = "terraform_tango_mysql"
    image = "ubuntu"
    flavor = "small"

    nics {
        network_id = "${var.network_id}"
    }
}