---
layout: "vra"
page_title: "VMware vRealize Automation: vra_image_profile"
description: |-
  Provides a data lookup for vra_image_profile.
---

# Resource: vra_image_profile
## Example Usages
This is an example of how to create an image profile resource.

**Image profile:**

```hcl
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

```

An image profile resource supports the following arguments:

## Argument Reference

* `description` - (Optional) A human-friendly description.

* `name` - (Required) A human-friendly name used as an identifier in APIs that support this option.

* `region_id` - (Required) The id of the region for which this profile is defined as in vRealize Automation(vRA).

## Attributes Reference

* `created_at` - Date when the entity was created. The date is in ISO 6801 and UTC.

* `external_region_id` - The external regionId of the resource. 

* `image_mapping` - Image mapping defined for the corresponding region.

* `owner` - Email of the user that owns the entity.

* `updated_at` - Date when the entity was last updated. The date is ISO 8601 and UTC.
