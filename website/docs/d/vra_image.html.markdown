---
layout: "vra"
page_title: "VMware vRealize Automation: vra_image"
description: |-
  Provides a data source for vRealize Automation images.
---

# Data Source: vra\_image

The `vra_image` data source can be used to discover the lookup machine images with cloud accounts. This can then be used with resource that require an image. For example, to create an image profile using the `vra_image_profile` resource.

## Example Usage

This is an example of how to lookup images.

**Image data source by Id:**

```hcl
# Lookup image using its id
data "vra_image" "this" {
  id = var.image_id
}
```

**Image data source by filter query:**

```hcl
# Lookup image using its name
data "vra_image" "this" {
  filter = "name eq '${var.image_name}'"
}
```

## Argument Reference

* `id` - (Optional) The id of the image resource instance. Only one of 'id' or 'filter' must be specified.
* `filter` - (Optional) Search criteria to narrow down the image resource instance. Only one of 'id' or 'filter' must be specified.

## Attribute Reference

* `created_at` - Date when the entity was created. The date is in ISO 8601 and UTC
* `custom_properties` - A list of key value pair of custom properties for the image resource.
* `description` - A human-friendly description.
* `external_id` - External entity Id on the provider side.
* `links` - HATEOAS of the entity.
* `name` - A human-friendly name used as an identifier for the image resource instance. 
* `org_id` - The id of the organization this entity belongs to.
* `os_family` - Operating System family of the image.
* `owner` - Email of the user that owns the entity.
* `private` - Indicates whether this image is private. For vSphere, private images are considered to be templates and snapshots and public are Content Library Items.
* `region` - The region of the image. For a vSphere cloud account, it is the `externalRegionId` such as `Datacenter:datacenter-2` and for an AWS cloud account, it is region name such as `us-east-1`, etc.
* `updated_at` - Date when the entity was last updated. The date is ISO 8601 and UTC.
