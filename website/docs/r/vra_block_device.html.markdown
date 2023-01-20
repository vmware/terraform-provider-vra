---
layout: "vra"
page_title: "VMware vRealize Automation: vra_block_device"
sidebar_current: "docs-vra-resource-block-device"
description: |-
  Creates a vra_block_device resource.
---

# Resource: vra_block_device

Creates a VMware vRealize Automation block device resource.

## Example Usages

The following example shows how to create a block device resource.

```hcl
resource "vra_block_device" "disk1" {
  capacity_in_gb = 10
  name = "terraform_vra_block_device1"
  project_id = var.project_id
  persistent = true
}
```

## Argument Reference

Create your block device resource with the following arguments:

* `capacity_in_gb` - (Required) Capacity of block device in GB.

* `name` - (Required) Human-friendly name used as an identifier in APIs that support this option.

* `project_id` - (Required) ID of project that current user belongs to.

* `constraints` - (Optional) Storage, network, and extensibility constraints to be applied when provisioning through the project.

* `description` - (Optional) Describes machine within the scope of your organization and is not propagated to the cloud.

* `disk_content_base_64` - (Optional) Content of a disk, base64 encoded.

* `encrypted` - (Optional) Indicates whether block device should be encrypted or not.

* `expand_snapshots` - (Optional) Indicates whether snapshots of block devices should be included in the state. Applies only to first class block devices.

* `purge` - (Optional) Indicates if the disk must be completely destroyed or should be kept in the system. Valid only for block devices with ‘persistent’ set to true. Used to destroy the resource.

* `persistent` - (Optional) Indicates whether block device survives a delete action.

* `source_reference` - (Optional) URI to use for block device. Example: `ami-0d4cfd66`

## Attribute Reference
* `created_at` - Date when entity was created. Date and time format is ISO 8601 and UTC.

* `custom_properties` - Additional custom properties that may be used to extend the machine.

* `deployment_id` - ID of deployment associated with resource.

* `external_id` - External entity ID on provider side.

* `external_region_id` - External regionId of resource.

* `external_zone_id` - External zoneId of resource.

* `org_id` - ID of organization that entity belongs to.

* `owner` - Email of entity owner.

* `links` - HATEOAS of entity.

* `snapshots` - Represents a machine snapshot.

    * `created_at` - Date when entity was created. Date and time format is ISO 8601 and UTC.

    * `description` - Describes machine within the scope of your organization and is not propagated to the cloud.

    * `id` - ID of the block device snapshot.

    * `is_current` - Indicates whether snapshot on block device is current.

    * `links` - HATEOAS of the entity

        * `rel`

        * `href`

        * `hrefs`

    * `name` - Human-friendly name for block device snapshot.

    * `org_id` - ID of organization that block device snapshot belongs to.

    * `owner` - Email of block device snapshot owner.

    * `updated_at` - Date when entity was last updated. Date and time format is ISO 8601 and UTC.

* `status` - Status of block device.

* `tags` - Set of tag keys and values to apply to the resource instance.
Example: `[ { "key" : "vmware.enumeration.type", "value": "nebs_block" } ]`

* `updated_at` - Date when entity was last updated. Date and time format is ISO 8601 and UTC.
