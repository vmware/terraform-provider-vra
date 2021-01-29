provider "vra" {
  url           = var.url
  refresh_token = var.refresh_token
  insecure      = var.insecure
}

// Required for vRA Cloud, Optional for vRA on-prem
data "vra_data_collector" "this" {
  count = var.data_collector_name != "" ? 1 : 0
  name  = var.data_collector_name
}

data "vra_region_enumeration_vmc" "this" {
  api_token = var.api_token
  sddc_name = var.sddc_name

  vcenter_hostname = var.vcenter_hostname
  vcenter_password = var.vcenter_password
  vcenter_username = var.vcenter_username
  dc_id            = var.data_collector_name != "" ? data.vra_data_collector.this[0].id : "" // Required for vRA Cloud, Optional for vRA on-prem
}

resource "vra_cloud_account_vmc" "this" {
  name        = "VMC cloud account"
  description = "Created via Terraform vRA provider"

  api_token = var.api_token
  sddc_name = var.sddc_name

  vcenter_hostname = var.vcenter_hostname
  vcenter_password = var.vcenter_password
  vcenter_username = var.vcenter_username
  nsx_hostname     = var.nsx_hostname
  dc_id            = "" // Required for vRA Cloud, Optional for vRA on-prem

  regions                 = data.vra_region_enumeration_vmc.this.regions
  accept_self_signed_cert = true
}

#Dummy wait for vsphere cloud account enumeration to complete
resource "null_resource" "delay" {
  provisioner "local-exec" {
    command = "sleep 60"
  }
  triggers = {
    "vc_account" = vra_cloud_account_vmc.this.id
  }
}

data "vra_region" "this" {
  // Important thing to note is the cloud_account_id in the filter criteria is ne (i.e. not equal) to VMC cloud account id
  filter = "externalRegionId eq '${var.external_region_id}' and cloudAccountId ne '${vra_cloud_account_vmc.this.id}'"
  name = var.region_name // Use 'name' argument in case there are other vCenter regions that have same external region id

  # filter by 'name' is not supported yet. Use the 'name' argument instead until it is supported
  # filter =  "name eq 'SDDC-Datacenter' and externalRegionId eq '${var.external_region_id}' and cloudAccountId ne vra_cloud_account_vmc.vmc_cloud_account.id"
}

resource "vra_zone" "this" {
  depends_on  = [null_resource.delay]
  name        = "VMC Zone"
  description = "Created via Terraform vRA provider"
  region_id   = data.vra_region.this.id
  folder      = "Workloads"

  tags {
    key   = "env"
    value = "prod"
  }
}

resource "vra_project" "this" {
  name        = "VMC Zone"
  description = "created via vra terraform provider"

  zone_assignments {
    zone_id       = vra_zone.this.id
    priority      = 1
    max_instances = 10
  }
  administrators = [var.user_email]
}

resource "vra_flavor_profile" "this" {
  name        = "VMC Flavor Profile"
  description = "Created via Terraform vRA provider"
  region_id   = data.vra_region.this.id

  flavor_mapping {
    name      = "small"
    cpu_count = "2"
    memory    = "2048"
  }
}

resource "vra_image_profile" "this" {
  name        = "VMC Image Profile"
  description = "Created via Terraform vRA provider"
  region_id   = data.vra_region.this.id

  image_mapping {
    name       = "ubuntu"
    image_name = "https://cloud-images.ubuntu.com/releases/16.04/release-20190605/ubuntu-16.04-server-cloudimg-amd64.ova"
  }
}

data "vra_fabric_network" "this" {
  depends_on = [vra_zone.this]
  filter     = "name eq '${var.fabric_network_name}'"
}

resource "vra_network_profile" "this" {
  depends_on  = [vra_zone.this]
  name        = "VMC Network Profile"
  description = "Created via Terraform vRA provider"
  region_id   = data.vra_region.this.id

  fabric_network_ids = [data.vra_fabric_network.this.id]
  isolation_type     = "NONE"

  tags {
    key   = "Dev"
    value = ""
  }
}

data "vra_fabric_datastore_vsphere" "this" {
  depends_on = [vra_zone.this]
  filter     = "name eq '${var.fabric_datastore_name}'"
}

resource "vra_storage_profile" "this" {
  depends_on          = [vra_zone.this]
  name                = "VMC Storage Profile"
  description         = "Created via Terraform vRA provider"
  region_id           = data.vra_region.this.id
  default_item        = true
  supports_encryption = true

  disk_properties = {
    independent      = true
    persistent       = true
    provisioningType = "thin"
    sharesLevel      = "High"
  }

  disk_target_properties = {
    datastoreId = data.vra_fabric_datastore_vsphere.this.id
  }
}