provider "vra" {
  url           = var.url
  refresh_token = var.refresh_token
  insecure      = var.insecure
}

resource "vra_project" "this" {
  name = var.project_name
}

resource "vra_content_source" "this" {
  name       = var.content_source_name
  project_id = vra_project.this.id
  // type_id needs to be one of com.gitlab, com.github or com.vmware.marketplace
  type_id     = "com.gitlab"
  description = "Some content Source"
  //whether automatically sync content or not
  sync_enabled = "false"
  config {
    path           = "blueprint01"
    branch         = "master"
    repository     = "vracontent/vra8_content_source_test"
    content_type   = "BLUEPRINT"
    project_name   = var.project_name
    integration_id = var.integration_id
  }


}
