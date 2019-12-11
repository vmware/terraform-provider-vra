---
layout: "vra"
page_title: "VMware vRealize Automation: vra_block_device"
sidebar_current: "docs-vra-resource-block-device"
description: |-
  Provides a VMware vRA vra_block_device resource.
---

# vra\_block\_device

Provides a VMware vRA vra_block_device resource.

## Example Usages

**Block device attached to a machine:**

This is an example on how to create a block device and attach it to a machine. The block devices and machine will be a part of the same deployment.

The region, flavors and image shown in the example are specific to AWS but these can be created for other cloud provider in a similar way. This example assumes that a cloud account, region (AWS in this case) and a project resource already exists. Look at the specific resource examples/documentation for more information.

```hcl

provider "vra" {
  url           = var.url
  refresh_token = var.refresh_token
}

data "vra_cloud_account_aws" "this" {
  name = var.cloud_account
}

data "vra_region" "this" {
  cloud_account_id = data.vra_cloud_account_aws.this.id
  region           = var.region
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

resource "vra_block_device" "disk1" {
  capacity_in_gb = 10
  name = "terraform_vra_block_device1"
  project_id = data.vra_project.this.id
}

resource "vra_block_device" "disk2" {
  capacity_in_gb = 10
  name = "terraform_vra_block_device2"
  project_id = data.vra_project.this.id
  deployment_id = vra_block_device.disk1.deployment_id
}

resource "vra_machine" "machine" {
  name        = "tf-machine"
  description = "terrafrom test machine"
  project_id  = data.vra_project.this.id
  image       = "ubuntu"
  flavor      = "small"
  deployment_id = vra_block_device.disk1.deployment_id

  tags {
    key   = "foo"
    value = "bar"
  }

  disks {
    block_device_id = vra_block_device.disk1.id
  }

   disks {
    block_device_id = vra_block_device.disk2.id
  }
}
```

## Argument Reference

The following arguments are supported for a machine resource:

* `capacity_in_gb` - (Required) Capacity of the block device in GB.

* `name` - (Required) A human-friendly name used as an identifier in APIs that support this option.

* `project_id` - (Required) The id of the project the current user belongs to.

* `constraints` - (Optional) Constraints that are used to drive placement policies for the virtual machine that is produced from this specification. Constraint expressions are matched against tags on existing placement targets.
example:[{"mandatory" : "true", "expression": "environment":"prod"}, {"mandatory" : "false", "expression": "pci"}]. It is nested argument with the following properties.
  * `expression` - (Required) A constraint that is conveyed to the policy engine. An expression of the form "[!]tag-key[:[tag-value]]", used to indicate a constraint match on keys and values of tags.
  * `mandatory` - (Required) Indicates whether this constraint should be strictly enforced or not.

* `custom_properties` - (Optional) Additional custom properties that may be used to extend the machine.
Following disk custom properties can be passed while creating a block device:
1. dataStore: Defines name of the datastore in which the disk has to be provisioned.
2. storagePolicy: Defines name of the storage policy in which the disk has to be provisioned. If name of the datastore is specified in the custom properties then, datastore takes precedence.
3. provisioningType: Defines the type of provisioning. For eg. thick/thin.

* `description` - (Optional) Describes machine within the scope of your organization and is not propagated to the cloud.

* `deployment_id` - (Optional) The id of the deployment that is associated with this resource.

* `disk_content_base_64` - (Optional) Content of a disk, base64 encoded.

* `encrypted` - (Optional) Indicates whether the block device should be encrypted or not.

* `source_reference` - (Optional) Reference to URI using which the block device has to be created.

* `tags` - (Optional) A set of tag keys and optional values that should be set on any resource that is produced from this specification. It is nested argument with the following properties.
  * `key` - (Required) Tag’s key.
  * `value` - (Required) Tag’s value.


## Attribute Reference

* `created_at` - Date when the entity was created. The date is in ISO 6801 and UTC.

* `external_id` - External entity Id on the provider side.

* `external_region_id` - The external regionId of the resource.

* `external_zone_id` - The external zoneId of the resource.

* `org_id` - The id of the organization this entity belongs to.

* `owner` - Email of the user that owns the entity.

* `links` - HATEOAS of the entity.

* `status` - Status of the block device.

* `updated_at` - Date when the entity was last updated. The date is ISO 8601 and UTC.


