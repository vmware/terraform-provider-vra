---layout: "vra"
page_title: "VMware vRealize Automation: vra_block_device"
description: |-
  Provides a data lookup for vra_block_device.
---

# Data Source: vra_block_device

Provides a data lookup for a vra_block_device.

## Example Usages

**Block device data source by its id:**

This is an example of how to read a block device data source using its ID.

```hcl
data "vra_block_device" "this" {
  id = var.block_device_id
}

```

**Block device data source filter by name:**

This is an example of how to read a block device data source using its name.

```hcl
data "vra_block_device" "this" {
  filter = "name eq '${var.block_device_name}'"
}

```

## Argument Reference

A block device data source supports the following arguments:

* `id` - (Optional) The id of the block device.

* `filter` - (Optional) Search criteria to filter the list of block devices.

* `expand_snapshots` - (Optional) Indicates whether the snapshots of the block-devices should be included in the state. Applicable only for first class block devices.

## Attributes Reference

A block device data source supports the following attributes:

* `capacity_in_gb` - Capacity of the block device in GB.

* `cloud_account_ids` - Set of ids of the cloud accounts this entity belongs to.

* `created_at` - Date when the entity was created. The date is in ISO 6801 and UTC.

* `custom_properties` - Additional custom properties that may be used to extend the machine.

* `deployment_id` - The id of the deployment that is associated with this resource.

* `description` - Describes machine within the scope of your organization and is not propagated to the cloud.

* `external_id` - External entity Id on the provider side.

* `external_region_id` - The external regionId of the resource.

* `external_zone_id` - The external zoneId of the resource.

* `name` - A human-friendly name used as an identifier in APIs that support this option.

* `org_id` - The id of the organization this entity belongs to.

* `owner` - Email of the user that owns the entity.

* `project_id` - The id of the project the current user belongs to.

* `persistent` - Indicates whether the block device survives a delete action.

* `links` - HATEOAS of the entity.

* `status` - Status of the block device.

* `tags` - A set of tag keys and optional values that were set on this resource instance.
example: `[ { "key" : "vmware.enumeration.type", "value": "nebs_block" } ]`
  * `key` - Tag’s key.
  * `value` - Tag’s value.

* `updated_at` - Date when the entity was last updated. The date is ISO 8601 and UTC.




