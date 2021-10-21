---
layout: "vra"
page_title: "VMware vRealize Automation: vra_block_device_snapshot"
description: |-
  Creates a VMware vRealize Automation vra_block_device_snapshot resource.
---

# Resource: vra_block_device_snapshot

Creates a VMware vRealize Automation block device snapshot resource.

## Example Usages

The following example shows how to create a block device snapshot resource.

```hcl
resource "vra_block_device_snapshot" "snapshot1" {
  block_device_id = var.block_device_id
  description = "terraform fcd snapshot"
}
```

## Argument Reference

Create your block device snapshot resource with the following arguments:

* `block_device_id` - (Required) ID of block device.

* `description` - (Optional) Human-friendly description.

## Attribute Reference
* `created_at` - Date when entity was created. Date and time format is ISO 8601 and UTC.

* `is_current` - Indicates whether snapshot on block device is current.

* `links` - HATEOAS of entity

* `name` - Human-friendly name used as an identifier in APIs that support this option.

* `org_id` - ID of organization that entity belongs to.

* `owner` - Email of entity owner.

* `update_at` - Date when entity was last updated. Date and time format is ISO 8601 and UTC.


