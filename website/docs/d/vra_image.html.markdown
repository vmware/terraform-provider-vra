---
layout: "vra"
page_title: "VMware vRealize Automation: vra_image"
description: |-
  Provides a data source for vRealize Automation images.
---

# Data Source: vra\_image

The `vra_image` data source can be used to discover the lookup machine images with cloud accounts. This can then be used with resource that require an image. For example, to create an image profile using the `vra_image_profile` resource.

## Example Usage
This is an example of how to lookup images from a vSphere cloud account.

```hcl
data "vra_cloud_account_vsphere" "this" {
  name = var.cloud_account
}

data "vra_region" "this" {
  cloud_account_id = data.vra_cloud_account_vsphere.this.id
  region = var.region
}

data "vra_image" "image_0" {
  filter = "name eq '${var.image_name_0}' and cloudAccountId eq '${data.vra_cloud_account_vsphere.this.id}' and externalRegionId eq '${var.region}'"
}

data "vra_image" "image_1" {
  filter = "name eq '${var.image_name_1}' and cloudAccountId eq '${data.vra_cloud_account_vsphere.this.id}' and externalRegionId eq '${var.region}'"
}

resource "vra_image_profile" "this" {
  name        = var.image_profile_name
  description = var.image_profile_description
  region_id   = data.vra_region.this.id

  image_mapping {
    name     = var.image_name_0
    image_id = data.vra_image.image_0.id
  }

  image_mapping {
    name     = var.image_name_1
    image_id = data.vra_image.image_1.id
  }
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Required) Search criteria to narrow down the image discovery.

## Attribute Reference

* `description` - A human-friendly description of the image.
* `external_id` - External entity id on the provider side.
* `id` - The id of the image.
* `name` - A human-friendly name used as an identifier in APIs that support this option.  
* `private` - Indicates whether this image is private. For vSphere, private images are templates and snapshots and public images are content library items.
* `region` - The regionId of the image. For a vSphere cloud account, it is the `externalRegionId` such as `Datacenter:datacenter-2` and for an AWS cloud account, it is region name such as `us-east-1`, etc.
