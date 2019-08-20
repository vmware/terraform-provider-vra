provider "vra" {
  url           = var.url
  refresh_token = var.refresh_token
}

resource "vra_cloud_account_aws" "my-cloud-account" {
  name        = "tf-vra-cloud-account-aws"
  description = "my terraform test cloud account aws"
  access_key  = var.access_key
  secret_key  = var.secret_key
  regions     = ["us-east-1", "us-west-1"]
}

data "vra_region" "us-east-1-region" {
  cloud_account_id = vra_cloud_account_aws.my-cloud-account.id
  region           = "us-east-1"
}

resource "vra_zone" "my-zone" {
  name        = "tf-vra-zone"
  description = "my terraform test cloud zone"
  region_id   = data.vra_region.us-east-1-region.id
  tags {
    key   = "my-tf-key"
    value = "my-tf-value"
  }
  tags {
    key   = "tf-foo"
    value = "tf-bar"
  }
}

resource "vra_project" "my-project" {
  name        = "tf-vra-project"
  description = "my terraform test project"
  zone_assignments {
    zone_id       = vra_zone.my-zone.id
    priority      = 1
    max_instances = 2
  }
}

resource "vra_flavor_profile" "my-flavor-profile" {
  name        = "tf-vra-flavor-profile"
  description = "my terraform test flavor profile"
  region_id   = data.vra_region.us-east-1-region.id
  flavor_mapping {
    name          = "small"
    instance_type = "t2.small"
  }
  flavor_mapping {
    name          = "medium"
    instance_type = "t2.medium"
  }
}

resource "vra_image_profile" "my-image-profile" {
  name        = "tf-vra-image-profile"
  description = "my terraform test image profile"
  region_id   = data.vra_region.us-east-1-region.id
  image_mapping {
    name       = "ubuntu"
    image_name = var.image_name
  }
}

resource "vra_block_device" "my-block-device" {
  name           = "tf-vra-block-device"
  capacity_in_gb = 4
  project_id     = vra_project.my-project.id

  tags {
    key   = "tf-foo"
    value = "tf-bar"
  }
}

