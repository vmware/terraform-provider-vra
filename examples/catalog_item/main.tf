provider "vra" {
  url           = var.url
  refresh_token = var.refresh_token
  insecure      = var.insecure
}

resource "random_integer" "suffix" {
  min = 1
  max = 50000
}

data "vra_zone" "this" {
  name = var.zone_name
}

resource "vra_project" "this" {
  name = format("%s-%d", var.project_name, random_integer.suffix.result)

  zone_assignments {
    zone_id          = data.vra_zone.this.id
    priority         = 1
    max_instances    = 2
    cpu_limit        = 1024
    memory_limit_mb  = 8192
    storage_limit_gb = 65536
  }
}

resource "vra_blueprint" "this" {
  name        = format("%s-%d", var.blueprint_name, random_integer.suffix.result)
  description = "Created by vRA terraform provider"

  project_id = vra_project.this.id

  content = <<-EOT
    formatVersion: 1
    inputs:
      image:
        type: string
        description: "Image"
      flavor:
        type: string
        description: "Flavor"
    resources:
      Machine:
        type: Cloud.Machine
        properties:
          image: $${input.image}
          flavor: $${input.flavor}
  EOT
}

// Example to create a blueprint version and release it
resource "vra_blueprint_version" "this" {
  blueprint_id = vra_blueprint.this.id
  description  = "Released from vRA terraform provider"
  version      = (random_integer.suffix.result / random_integer.suffix.result)
  release      = true
  change_log   = "First version"
}

// Example to fetch a blueprint version and release it
data "vra_blueprint_version" "this" {
  blueprint_id = vra_blueprint.this.id
  id           = vra_blueprint_version.this.id
}

// Example to create a blueprint content source for a project in Service Broker
resource "vra_catalog_source_blueprint" "this" {
  depends_on = [vra_blueprint_version.this]
  name       = format("%s-%d", var.catalog_source_name, random_integer.suffix.result)
  project_id = vra_project.this.id
}

// Example to fetch a blueprint catalog soruce id by project_id
data "vra_catalog_source_blueprint" "this" {
  depends_on = [vra_catalog_source_blueprint.this]
  project_id = vra_catalog_source_blueprint.this.project_id
}

// Example to create a content sharing for a blueprint content source
resource "vra_content_sharing_policy" "catalog_source_entitlement" {
  depends_on         = [vra_catalog_source_blueprint.this]
  name               = vra_catalog_source_blueprint.this.name
  project_id         = vra_project.this.id
  catalog_source_ids = [vra_catalog_source_blueprint.this.id]
}

// Example to fetch a content sharing/entitlement by id
data "vra_content_sharing_policy" "this" {
  id = vra_content_sharing_policy.catalog_source_entitlement.id
}

// Example to fetch a catalog item
data "vra_catalog_item" "this" {
  depends_on      = [vra_catalog_source_blueprint.this, vra_catalog_source_entitlement.this]
  name            = vra_blueprint.this.name
  expand_versions = true
}

// Example to create a content sharing for a content item
resource "vra_content_sharing_policy" "catalog_item_entitlement" {
  name             = vra_catalog_item.this.name
  project_id       = vra_project.this.id
  catalog_item_ids = [data.vra_catalog_item.this.id]
}

// Example to fetch a content sharing/entitlement by name
data "vra_content_sharing_policy" "this" {
  name = vra_content_sharing_policy.catalog_item_entitlement.name
}

// Example to request a deployment from a catalog item
resource "vra_deployment" "this" {
  name        = format("%s-%d", var.deployment_name, random_integer.suffix.result)
  description = "terraform test deployment"

  catalog_item_id      = data.vra_catalog_item.this.id
  catalog_item_version = vra_blueprint_version.this.version
  project_id           = vra_project.this.id

  inputs = {
    flavor = "small"
    image  = "centos"
  }

  timeouts {
    create = "30m"
    delete = "30m"
    update = "30m"
  }
}