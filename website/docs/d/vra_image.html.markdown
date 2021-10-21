---
layout: "vra"
page_title: "VMware vRealize Automation: vra_image"
description: |-
  Provides a data lookup for vRA Images.
---

# Data Source: vra_image
## Example Usages
This is an example of how to lookup Images.

**Images by Id:**

```hcl
# Lookup Images using its Id
data "vra_image" "this" {
  id = var.image_id
}
```

**Images by filter query:**

```hcl
# Lookup Images using its name
data "vra_image" "this" {
  filter = "name eq '${var.name}'"
}
```

An Image supports the following arguments:

## Argument Reference
* `filter` - Search criteria to narrow down Images.

* `id` - The id of the Image.

## Attribute Reference

* `description` - A human-friendly description of the fabric vSphere storage account.

* `external_id` - External entity Id on the provider side.

* `name` - A human-friendly name used as an identifier in APIs that support this option.  

* `private` - Indicates whether this fabric image is private. For vSphere, private images are considered to be templates and snapshots and public are Content Library Items.

* `region` - The regionId of the image.