---
layout: "vra"
page_title: "VMware vRealize Automation: vra_machine"
sidebar_current: "docs-vra-resource-machine"
description: |-
  Provides a VMware vRA vra_machine resource.
---

# vra\_machine

Provides a VMware vRA vra_machine resource.

## Example Usages

**Simple cloud-agnostic machine:**

This is an example on how to create a cloud agnostic machine along with an image and a flavor profile.
Image profile represents a structure that holds a list of image mappings defined for the particular region.
Flavor profile represents a structure that holds flavor mappings defined for the corresponding cloud end-point region.

```hcl
data "vra_cloud_account_aws" "this" {
  name = var.cloud_account
}

data "vra_region" "this" {
  cloud_account_id = data.vra_cloud_account_aws.this.id
  region           = var.region
}

data "vra_zone" "this" {
  name = var.zone
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

resource "vra_machine" "this" {
  name        = "tf-machine"
  description = "terrafrom test machine"
  project_id  = data.vra_project.this.id
  image       = "ubuntu"
  flavor      = "small"

  tags {
    key   = "foo"
    value = "bar"
  }
}
```

## Argument Reference

The following arguments are supported for a machine resource:

* `flavor` - (Required) Flavor of machine instance

* `image` - (Required) The type of image used for this machine.

* `name` - (Required) A human-friendly name used as an identifier in APIs that support this option.

* `project_id` - (Required) The id of the project the current user belongs to.

* `boot_config` - (Optional) Machine boot config that will be passed to the instance that can be used to perform common automated configuration tasks and even run scripts after the instance starts. It is nested argument with the following property.
  * `content` - A valid cloud config data in json-escaped yaml syntax

* `constraints` - (Optional) Constraints that are used to drive placement policies for the virtual machine that is produced from this specification. Constraint expressions are matched against tags on existing placement targets.
example:[{"mandatory" : "true", "expression": "environment":"prod"}, {"mandatory" : "false", "expression": "pci"}]. It is nested argument with the following properties.
  * `expression` - A constraint that is conveyed to the policy engine. An expression of the form "[!]tag-key[:[tag-value]]", used to indicate a constraint match on keys and values of tags. 
  * `mandatory` - Indicates whether this constraint should be strictly enforced or not.

* `custom_properties` - (Optional) Additional custom properties that may be used to extend the machine.

* `description` - (Optional) Describes machine within the scope of your organization and is not propagated to the cloud.

* `disks` - (Optional) A set of disk specifications for this machine. It is nested argument with the following properties.
  * `block_device_id` - The id of the existing block device.
  * `description` - A human-friendly description.
  * `name` - A human-friendly name used as an identifier in APIs that support this option.

* `image_disk_constraints` - (Optional) Constraints that are used to drive placement policies for the image disk. Constraint expressions are matched against tags on existing placement targets. example:[{"mandatory" : "true", "expression": "environment:prod"}, {"mandatory" : "false", "expression": "pci"}]. It is nested argument with the following properties.
  * `expression` - A constraint that is conveyed to the policy engine. An expression of the form "[!]tag-key[:[tag-value]]", used to indicate a constraint match on keys and values of tags. 
  * `mandatory` - Indicates whether this constraint should be strictly enforced or not.

* `image_ref` - (Optional) Direct image reference used for this machine (name, path, location, uri, etc.). Valid if no image property is provided

* `nics` - (Optional) Specification for attaching nic to machine. A set of network interface controller specifications for this machine. If not specified, then a default network connection will be created. It is nested argument with the following properties.
  * `addresses` - A list of IP addresses allocated or in use by this network interface.
  * `customProperties` - Additional properties that may be used to extend the base type.
  * `securityGroupIds` - list of security group ids which this network interface will be assigned to.
  * `name` - A human-friendly name used as an identifier in APIs that support this option.
  * `description` - A human-friendly description.
  * `networkId` - Id of the network instance that this network interface plugs into.
  * `deviceIndex` - The device index of this network interface.

* `tags` - (Optional) A set of tag keys and optional values that should be set on any resource that is produced from this specification.

* `address` - (Computed) Primary address allocated or in use by this machine. The actual type of the address depends on the adapter type. Typically it is either the public or the external IP address.

* `created_at` - (Computed) Date when the entity was created. The date is in ISO 6801 and UTC.

* `external_id` - (Computed) External entity Id on the provider side.

* `external_region_id` - (Computed) The external regionId of the resource.

* `external_zone_id` - (Computed) The external zoneId of the resource.

* `organization_id` - (Computed) The id of the organization this entity belongs to.

* `owner` - (Computed) Email of the user that owns the entity.

* `power_state` - (Computed) Power state of machine.

* `updated_at` - (Computed) Date when the entity was last updated. The date is ISO 8601 and UTC.

