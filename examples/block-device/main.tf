provider "vra" {
  url           = var.url
  refresh_token = var.refresh_token
}

data "vra_cloud_account_aws" "this" {
  name = var.cloud_account
}

data "vra_region" "this" {
  cloud_account_id = data.vra_cloud_account_aws.this.id
  region           = var.region
}

data "vra_project" "this" {
  name = var.project
}

resource "vra_flavor_profile" "this" {
  name        = "tf-vra-flavor-profile"
  description = "my flavor"
  region_id   = data.vra_region.this.id

  flavor_mapping {
    name          = "small"
    instance_type = "t2.small"
  }

  flavor_mapping {
    name          = "medium"
    instance_type = "t2.medium"
  }
}

resource "vra_image_profile" "this" {
  name        = "tf-vra-image-profile"
  description = "terraform test image profile"
  region_id   = data.vra_region.this.id

  image_mapping {
    name       = "ubuntu"
    image_name = var.image_name
  }
}

resource "vra_block_device" "disk1" {
  capacity_in_gb = 10
  name = "terraform_vra_block_device1"
  project_id = data.vra_project.this.id
  persistent = true
}

resource "vra_block_device" "disk2" {
  capacity_in_gb = 10
  name = "terraform_vra_block_device2"
  project_id = data.vra_project.this.id
  deployment_id = vra_block_device.disk1.deployment_id
}

resource "vra_machine" "machine" {
  name        = "tf-machine"
  description = "terraform test machine"
  project_id  = data.vra_project.this.id
  image       = "ubuntu"
  flavor      = "small"
  deployment_id = vra_block_device.disk1.deployment_id

  tags {
    key   = "foo"
    value = "bar"
  }

  disks {
    block_device_id = vra_block_device.disk1.id
  }

   disks {
    block_device_id = vra_block_device.disk2.id
  }
}