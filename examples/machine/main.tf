provider "vra" {
  url           = var.url
  refresh_token = var.refresh_token
}

data "vra_cloud_account_aws" my-cloud-account {
       name = var.cloud_account
}

data "vra_region" "us-east-1-region" {
  cloud_account_id = vra_cloud_account_aws.my-cloud-account.id
  region           = var.region
}

data "vra_zone" "my-zone" {
  name = var.zone
}

data "vra_project" "my-project" {
  name = var.project
}

resource "vra_flavor_profile" "my-flavor-profile" {
	name        = "tf-vra-flavor-profile"
	description = "my flavor"
	region_id   = "${data.vra_region.us-east-1-region.id}"

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

resource "vra_machine" "my-machine" {
    name        = "tf-machine"
    description = "my terrafrom test machine"
    project_id  = data.vra_project.my-project.id
    image       = "ubuntu"
    flavor      = "small"

    tags {
      key   = "foo"
      value = "bar"
    }
}
