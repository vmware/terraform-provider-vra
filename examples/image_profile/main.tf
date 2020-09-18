provider "vra" {
  url           = var.url
  refresh_token = var.refresh_token
}

data "vra_cloud_account_vsphere" "this" {
  name = var.cloud_account
}

data "vra_region" "this" {
  cloud_account_id = data.vra_cloud_account_vsphere.this.id
  region           = var.region
}

data "vra_image" "centos" {
  filter = "name eq '${var.image_name1}' and cloudAccountId eq '${data.vra_cloud_account_vsphere.this.id}' and externalRegionId eq '${var.region}'"
}

data "vra_image" "photon" {
  filter = "name eq '${var.image_name2}' and cloudAccountId eq '${data.vra_cloud_account_vsphere.this.id}'"
}

resource "vra_image_profile" "this" {
  name        = "vra-image-profile"
  description = "test image profile"
  region_id   = data.vra_region.this.id

  image_mapping {
    name     = "centos"
    image_id = data.vra_image.centos.id

    constraints {
      mandatory  = true
      expression = "!env:Test"
    }
    constraints {
      mandatory  = false
      expression = "foo:bar"
    }
  }

  image_mapping {
    name     = "photon"
    image_id = data.vra_image.photon.id

    cloud_config = "runcmd echo 'Hello'"
  }
}

// Image profile data source by region id
data "vra_image_profile" "this" {
  region_id = vra_image_profile.this.region_id

  depends_on = [vra_image_profile.this]
}

// Image profile data source by name
data "vra_image_profile" "name" {
  name = vra_image_profile.this.name

  depends_on = [vra_image_profile.this]
}

// Image profile data source by filter
data "vra_image_profile" "filter" {
  filter = "regionId eq '${vra_image_profile.this.region_id}'"

  depends_on = [vra_image_profile.this]
}

// Image profile data source by id
data "vra_image_profile" "id" {
  id = vra_image_profile.this.id
}