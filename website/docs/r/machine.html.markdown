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

The region, flavors and image shown in the example are specific to AWS but these can be cretaed for other cloud provider in a similar way. This example assumes that a cloud account, zone (AWS in this case) and a project resource already exists. Look at the specific resource examples/documentation for more information.

```hcl

data "vra_cloud_account_aws" "this" {
  name = "example-cloud-account-aws"
}

data "vra_region" "this" {
  cloud_account_id = data.vra_cloud_account_aws.this.id
  region           = "us-east-1"
}

data "vra_project" "this" {
  name = "example-project"
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
    name       = "centos"
    image_name = "ami-e416a39e"
  }
}

data "vra_network" "this" {
  name = "appnet-public"
}

resource "vra_machine" "this" {
  name        = "tf-vra-machine"
  description = "terrafrom test machine"
  project_id  = data.vra_project.this.id
  image       = "centos"
  flavor      = "small"

  nics {
      network_id = data.vra_network.this.id
  }

  boot_config {
      content = <<EOF
  #cloud-config
    users:
    - default
    - name: myuser
      sudo: ['ALL=(ALL) NOPASSWD:ALL']
      groups: [wheel, sudo, admin]
      shell: '/bin/bash'
      ssh-authorized-keys: |
        ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDytVL+Q6/UmGwdnJxKQEozqERHGqxlH/zBbT5W8iNbwgOLF6JWz0o7ThAK/Cf0uPcv78Q6UhOjuRfd2BKBciJx5JsyH4Ly7Ars2v/ZQ492KyZElKRqwibXNWjfZcwKU/6YjDITm15Yh6UWCsvVHg4w72X+TiTxeKDZ0pNt2hcZ5Uje6NvZ4GFKYfl4kNFxBZmBYLFdtq8eNPg3PGREV+pM0xkyXKSAYUsXsgj821AgK/YNByCPY53jNKqXqdFKQXKG7FOs78MdhAF7aGMsVRymY5RtHk9UO0DGzCIHRp9DqmfN9SdIYIf5fb4sEtt8T9uxW32Mx3d9S+vGbmkYoRpY user@example.com
  
    runcmd:
      - sudo sed -e 's/.*PasswordAuthentication yes.*/PasswordAuthentication no/' -i /etc/ssh/sshd_config
      - sudo service sshd restart
  EOF
    }
  
    constraints {
      mandatory  = true
      expression = "AWS"
    }
  
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
  * `content` - (Optional) A valid cloud config data in json-escaped yaml syntax

* `constraints` - (Optional) Constraints that are used to drive placement policies for the virtual machine that is produced from this specification. Constraint expressions are matched against tags on existing placement targets.
example:[{"mandatory" : "true", "expression": "environment":"prod"}, {"mandatory" : "false", "expression": "pci"}]. It is nested argument with the following properties.
  * `expression` - (Required) A constraint that is conveyed to the policy engine. An expression of the form "[!]tag-key[:[tag-value]]", used to indicate a constraint match on keys and values of tags.
  * `mandatory` - (Required) Indicates whether this constraint should be strictly enforced or not.

* `custom_properties` - (Optional) Additional custom properties that may be used to extend the machine.

* `description` - (Optional) Describes machine within the scope of your organization and is not propagated to the cloud.

* `deployment_id` - (Optional) The id of the deployment that is associated with this resource.

* `disks` - (Optional) A set of disk specifications for this machine. It is nested argument with the following properties.
  * `block_device_id` - (Required) The id of the existing block device.
  * `description` - (Optional) A human-friendly description.
  * `name` - (Optional) A human-friendly name used as an identifier in APIs that support this option.

* `image_disk_constraints` - (Optional) Constraints that are used to drive placement policies for the image disk. Constraint expressions are matched against tags on existing placement targets. example:[{"mandatory" : "true", "expression": "environment:prod"}, {"mandatory" : "false", "expression": "pci"}]. It is nested argument with the following properties.
  * `expression` - (Required) A constraint that is conveyed to the policy engine. An expression of the form "[!]tag-key[:[tag-value]]", used to indicate a constraint match on keys and values of tags.
  * `mandatory` - (Required) Indicates whether this constraint should be strictly enforced or not.

* `image_ref` - (Optional) Direct image reference used for this machine (name, path, location, uri, etc.). Valid if no image property is provided

* `nics` - (Optional) Specification for attaching nic to machine. A set of network interface controller specifications for this machine. If not specified, then a default network connection will be created. It is nested argument with the following properties.
  * `addresses` - (Optional) A list of IP addresses allocated or in use by this network interface.
  * `customProperties` - (Optional) Additional properties that may be used to extend the base type.
  * `securityGroupIds` - (Optional) List of security group ids which this network interface will be assigned to.
  * `name` - (Optional) A human-friendly name used as an identifier in APIs that support this option.
  * `description` - (Optional) A human-friendly description.
  * `networkId` - (Required) Id of the network instance that this network interface plugs into.
  * `deviceIndex` - (Optional) The device index of this network interface.

* `tags` - (Optional) A set of tag keys and optional values that should be set on any resource that is produced from this specification. It is nested argument with the following properties.
  * `key` - (Required) Tag’s key.
  * `value` - (Required) Tag’s value.


## Attribute Reference

* `address` - Primary address allocated or in use by this machine. The actual type of the address depends on the adapter type. Typically it is either the public or the external IP address.

* `created_at` - Date when the entity was created. The date is in ISO 6801 and UTC.

* `external_id` - External entity Id on the provider side.

* `external_region_id` - The external regionId of the resource.

* `external_zone_id` - The external zoneId of the resource.

* `organization_id` - The id of the organization this entity belongs to.

* `owner` - Email of the user that owns the entity.

* `power_state` - Power state of machine.

* `updated_at` - Date when the entity was last updated. The date is ISO 8601 and UTC.

