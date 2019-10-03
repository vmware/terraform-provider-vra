---
layout: "vra"
page_title: "VMware vRealize Automation: vra_machine"
sidebar_current: "docs-vra-resource-machine"
description: |-
  Provides a VMware vRA vra_machine resource.
---

# vra\_machine

Provides a VMware vRA vra_machine resource.

## Example Usages

**Simple cloud-agnostic machine:**

This is an example on how to create a cloud agnostic machine along with an image and a flavor profile.
Image profile represents a structure that holds a list of image mappings defined for the particular region.
Flavor profile represents a structure that holds flavor mappings defined for the corresponding cloud end-point region.

```hcl
data "vra_cloud_account_aws" "this" {
  name = var.cloud_account
}

data "vra_region" "this" {
  cloud_account_id = data.vra_cloud_account_aws.this.id
  region           = var.region
}

data "vra_zone" "this" {
  name = var.zone
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

resource "vra_machine" "this" {
  name        = "tf-machine"
  description = "terrafrom test machine"
  project_id  = data.vra_project.this.id
  image       = "ubuntu"
  flavor      = "small"

  tags {
    key   = "foo"
    value = "bar"
  }
}
```

## Argument Reference

The following arguments are supported for a machine resource:

* `name` - (Required) A human-friendly name used as an identifier in APIs that support this option.
* `description` - (Optional) Describes machine within the scope of your organization and is not propagated to the cloud.
* `project_id` - (Required) The id of the project the current user belongs to.
* `image` - (Required) The type of image used for this machine.
* `flavor` - (Required) Flavor of machine instance
* `tags` - (Optional) A set of tag keys and optional values that should be set on any resource that is produced from this specification.
