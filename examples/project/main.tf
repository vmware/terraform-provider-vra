provider "vra" {
  url           = var.url
  refresh_token = var.refresh_token
  insecure      = var.insecure
}

data "vra_zone" "this" {
  name = var.zone_name
}

resource "vra_project" "this" {
  name        = var.project_name
  description = "terraform test project"

  zone_assignments {
    zone_id          = data.vra_zone.this.id
    priority         = 1
    max_instances    = 2
    cpu_limit        = 1024
    memory_limit_mb  = 8192
    storage_limit_gb = 65536
  }

  shared_resources = true

  administrators = ["jason@vra.local"]
}
