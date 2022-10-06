provider "vra" {
  url           = var.url
  refresh_token = var.refresh_token
}

data "vra_cloud_account_aws" "aws" {
  name = var.aws_cloud_account
}

data "vra_cloud_account_vsphere" "vsphere" {
  name = var.vsphere_cloud_account
}

data "vra_region" "aws" {
  cloud_account_id = data.vra_cloud_account_aws.aws.id
  region           = var.aws_region
}

data "vra_region" "vsphere" {
  cloud_account_id = data.vra_cloud_account_vsphere.vsphere.id
  region           = var.vsphere_region
}

resource "vra_flavor_profile" "aws" {
  name        = "tf-aws-flavor-profile"
  description = "My AWS flavor"
  region_id   = data.vra_region.aws.id

  flavor_mapping {
    name          = "small"
    instance_type = "t2.small"
  }

  flavor_mapping {
    name          = "medium"
    instance_type = "t2.medium"
  }
}

resource "vra_flavor_profile" "vsphere" {
  name        = "tf-vsphere-flavor-profile"
  description = "My vSphere flavor"
  region_id   = data.vra_region.vsphere.id

  flavor_mapping {
    name      = "small"
    cpu_count = 2
    memory    = 2048
  }

  flavor_mapping {
    name      = "medium"
    cpu_count = 4
    memory    = 4096
  }
}
