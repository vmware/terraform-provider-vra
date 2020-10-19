---
layout: "vra"
page_title: "VMware vRealize Automation: vra_block_device"
description: |-
  Provides a VMware vRA vra_block_device resource.
---

# Resource: vra_block_device_snapshots
## Example Usages

This is an example of how to create a block device snapshot resource.

```hcl
resource "vra_block_device_snapshot" "snapshot1" {
  block_device_id = vra_block_device.disk1.id
  description = "terraform fcd snapshot"
}
```
A block device snapshot resource supports the following resource:

## Argument Reference

* `block_device_id` - The id of the block device.

* `description` - A human-friendly description.

## Attribute Reference
* `created_at` - Date when the entity was created. The date is in ISO 6801 and UTC.

* `is_current` - Indicates whether this snapshot is the current snapshot on the block-device.

* `links` - HATEOAS of the entity

* `name` - A human-friendly name used as an identifier in APIs that support this option.  Only one of 'filter', 'id', 'name' or 'region_id' must be specified.

* `org_id` - The id of the organization this entity belongs to.

* `owner` - Email of the user that owns the entity.

* `update_at` - Date when the entity was last updated. The date is ISO 8601 and UTC.


