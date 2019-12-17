provider "vra" {
  url           = var.url
  refresh_token = var.refresh_token
}

data "vra_project" "this" {
  name = var.project_name
}

data "vra_blueprint" "this" {
  name = var.blueprint_name
}

resource "vra_deployment" "this" {
  name        = var.deployment_name
  description = "Deployed from vRA provider for Terraform."

  blueprint_id      = data.vra_blueprint.this.id
  blueprint_version = var.blueprint_version
  project_id        = data.vra_project.this.id

  inputs = {
    flavor = "small"
    image  = "centos"
    count  = 2
    flag   = true
  }

  expand_resources    = true
  expand_last_request = true
}
