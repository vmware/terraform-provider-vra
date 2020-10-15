---layout: "vra"
page_title: "VMware vRealize Automation: vra_block_device"
sidebar_current: "docs-vra-datasource-vra_block_device"
description: |-
  Provides a data lookup for vra_block_device.
---

# Data Source: vra\_block\_device

Provides a VMware vRA vra_block_device data source.

## Example Usages

**Block device data source by its id:**

This is an example of how to create a block device resource and read it as a data source using its id.
NOTE: The block device resource need not be created through terraform.
To create a block device, follow the resource block device documentation

```hcl

provider "vra" {
  url           = var.url
  refresh_token = var.refresh_token
}

data "vra_project" "this" {
  name = var.project
}

resource "vra_block_device" "this" {
  capacity_in_gb = 10
  name = "terraform_vra_block_device1"
  project_id = data.vra_project.this.id
}

data "vra_block_device" "this" {
  id = "vra_block_device.id"
}

```

**Block device data source filter by name:**

This is an example of how to create a block device resource and filter it as a data source using its name.
NOTE: The block device resource need not be created through terraform.
To create a block device, follow the resource block device documentation.

```hcl

provider "vra" {
  url           = var.url
  refresh_token = var.refresh_token
}

data "vra_project" "this" {
  name = var.project
}

resource "vra_block_device" "this" {
  capacity_in_gb = 10
  name = "terraform_vra_block_device1"
  project_id = data.vra_project.this.id
}

data "vra_block_device" "this" {
  filter = "name eq '${vra_block_device.this.name}'"
}

```



## Argument Reference

The following arguments are supported for a machine resource:

* `id` - (Optional) The id of this resource instance.

* `filter` - (Optional) The filter query


## Attribute Reference

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
example:[ { "key" : "vmware.enumeration.type", "value": "nebs_block" } ]
  * `key` - Tag’s key.
  * `value` - Tag’s value.

* `updated_at` - Date when the entity was last updated. The date is ISO 8601 and UTC.




