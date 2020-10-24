---
layout: "vra"
page_title: "VMware vRealize Automation: vra_block_device"
sidebar_current: "docs-vra-resource-block-device"
description: |-
  Provides a VMware vRA vra_block_device resource.
---

# vra_block_device
## Example Usages

# Resource: vra_block_device
## Example Usages

This is an example of how to read a block device resource.

```hcl
resource "vra_block_device" "disk1" {
  capacity_in_gb = 10
  name = "terraform_vra_block_device1"
  project_id = data.vra_project.this.id
  persistent = true
}
```

A block device resource supports the following arguments:

## Argument Reference
* `capacity_in_gb` - (Required) Capacity of the block device in GB.

* `name` - (Required) A human-friendly name used as an identifier in APIs that support this option.

* `project_id` - (Required) The id of the project the current user belongs to.

* `constraints` - (Optional) List of storage, network and extensibility constraints to be applied when provisioning through this project.

* `description` - (Optional) Describes machine within the scope of your organization and is not propagated to the cloud.

* `disk_content_base_64` - (Optional) Content of a disk, base64 encoded.

* `encrypted` - (Optional) Indicates whether the block device should be encrypted or not.

* `expand_snapshots` - (Optional) Indicates whether the snapshots of the block-devices should be included in the state. Applicable only for first class block devices.

* `purge` - (Optional) Indicates if the disk has to be completely destroyed or should be kept in the system. Valid only for block devices with ‘persistent’ set to true, only used for destroy the resource.

* `persistent` - (Optional) Indicates whether the block device survives a delete action.

* `source_reference` - (Optional) Reference to URI using which the block device has to be created. Example: ami-0d4cfd66

## Attribute Reference
* `created_at` - Date when the entity was created. The date is in ISO 6801 and UTC.

* `custom_properties` - Additional custom properties that may be used to extend the machine.

* `deployment_id` - The id of the deployment that is associated with this resource.

* `external_id` - External entity Id on the provider side.

* `external_region_id` - The external regionId of the resource.

* `external_zone_id` - The external zoneId of the resource.

* `org_id` - The id of the organization this entity belongs to.

* `owner` - Email of the user that owns the entity.

* `links` - HATEOAS of the entity.

* `snapshots` - Represents a machine snapshot.
    
    * `created_at` - Date when the entity was created. The date is in ISO 6801 and UTC.
    
    * `description` - Describes machine within the scope of your organization and is not propagated to the cloud.
    
    * `id` - The id of this block-device snapshot.
    
    * `is_current` - Indicates whether this snapshot is the current snapshot on the block-device.
    
    * `links` - HATEOAS of the entity
    
        * `rel`
        
        * `href`
        
        * `hrefs`
    
    * `name` - A human-friendly name for the block-device snapshot.
    
    * `org_id` - The id of the organization this block-device snapshot belongs to.
    
    * `owner` - Email of the user user that owns the block-device snapshot.
    
    * `updated_at` - Date when the entity was last updated. The date is in ISO 8601 and UTC.

* `status` - Status of the block device.

* `tags` - A set of tag keys and optional values that were set on this resource instance.
example:[ { "key" : "vmware.enumeration.type", "value": "nebs_block" } ]
  * `key` - Tag’s key.
  * `value` - Tag’s value.

* `updated_at` - Date when the entity was last updated. The date is ISO 8601 and UTC.
