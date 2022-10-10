provider "vra" {
  url           = var.url
  refresh_token = var.refresh_token
  insecure      = var.insecure
}

resource "vra_integration" "ad" {
  name        = "active-directory"
  description = "Active Directory Integration"
  integration_type = "activedirectory"
  integration_properties = {
    server: vra.ad_server
    endpointId: vra.ad_endpoint_id
    user: var.ad_user
    defaultOU: var.ad_default_ou
  }
  private_key = var.ad_password
}

resource "vra_integration" "github" {
  name        = "github"
  description = "Github Integration"
  integration_type = "com.github.saas"
  integration_properties = {
    url: "https://api.github.com"
  }
  private_key = var.github_token
}

resource "vra_integration" "saltstack" {
  name        = "saltstack"
  description = "SaltStack Integration"
  integration_type = "saltstack"
  integration_properties = {
    hostName: vra.saltstack_hostname
    endpointId: vra.saltstack_endpoint_id
  }
  private_key_id = var.saltstack_username
  private_key = var.saltstack_password
}

