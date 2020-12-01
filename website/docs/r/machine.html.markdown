---
layout: "vra"
page_title: "VMware vRealize Automation: vra_machine"
sidebar_current: "docs-vra-resource-machine"
description: |-
  Creates a vra_machine resource.
---

# vra\_machine

Creates a VMware vRealize Automation machine resource.

## Example Usages

**Simple cloud-agnostic machine:**

The following example shows how to create a cloud agnostic machine along with an image and a flavor profile.
* Image profile represents a structure that holds a list of image mappings defined for the particular region.
* Flavor profile represents a structure that holds flavor mappings defined for the corresponding cloud endpoint region.

The region, flavors, and image shown in the example are specific to AWS but you can create machines for other cloud providers in a similar way. This example assumes that a cloud account, cloud zone, and project resource already exist. See specific resource examples or documentation for more information.

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

Create your machine resource with the following arguments:

* `flavor` - (Required) Flavor of machine instance.

* `image` - (Required) Type of image used for the machine.

* `name` - (Required) Human-friendly name used as an identifier in APIs that support this option.

* `project_id` - (Required) ID of project that current user belongs to.

* `boot_config` - (Optional) Machine boot config to be passed to the instance. Used to perform common automated configuration tasks or run scripts after the instance starts. Nested argument with the following property:
  * `content` - (Optional) Valid cloud config data in JSON-escaped YAML syntax

* `constraints` - (Optional) Constraints used to drive placement policies for the virtual machine produced from the specification. Constraint expressions are matched against tags on existing placement targets.  
Example:[{"mandatory" : "true", "expression": "environment":"prod"}, {"mandatory" : "false", "expression": "pci"}].  
Nested argument with the following properties.
  * `expression` - (Required) Constraint conveyed to the policy engine. Expression of the form "[!]tag-key[:[tag-value]]", indicates a constraint match on keys and values of tags.
  * `mandatory` - (Required) Indicates whether constraint must be strictly enforced.

* `custom_properties` - (Optional) Additional custom properties that may be used to extend the machine.

* `description` - (Optional) Describes machine within the scope of your organization and is not propagated to the cloud.

* `deployment_id` - (Optional) ID of the deployment associated with the resource.

* `disks` - (Optional) Set of disk specifications for the machine. Nested argument with the following properties:
  * `block_device_id` - (Required) ID of the existing block device.
  * `description` - (Optional) Human-friendly description.
  * `name` - (Optional) Human-friendly name used as an identifier in APIs that support this option.

* `image_disk_constraints` - (Optional) Constraints used to drive placement policies for the image disk. Constraint expressions are matched against tags on existing placement targets.  
Example:[{"mandatory" : "true", "expression": "environment:prod"}, {"mandatory" : "false", "expression": "pci"}].  
Nested argument with the following properties.
  * `expression` - (Required) Constraint conveyed to the policy engine. Expression of the form "[!]tag-key[:[tag-value]]", indicates a constraint match on keys and values of tags.
  * `mandatory` - (Required) Indicates whether constraint must be strictly enforced.

* `image_ref` - (Optional) Direct image reference used for this machine such as name, path, location, URI. Valid if no image property is provided

* `nics` - (Optional) Set of network interface controller specifications used to attach a NIC to the machine. If left unspecified, a default network connection is created. Nested argument with the following properties.
  * `addresses` - (Optional) IP addresses allocated or in use by the network interface.
  * `customProperties` - (Optional) Additional properties used to extend the base type.
  * `securityGroupIds` - (Optional) Security group IDs that the network interface will be assigned to.
  * `name` - (Optional) Human-friendly name used as an identifier in APIs that support this option.
  * `description` - (Optional) Human-friendly description.
  * `networkId` - (Required) ID of the network instance that the network interface plugs into.
  * `deviceIndex` - (Optional) Device index of the network interface.

* `tags` - (Optional) Set of tag keys and values to apply to any resource produced from the specification.  
Example:[ { "key" : "foo", "value": "bar" } ]


## Attribute Reference

* `address` - Primary address allocated or in use by the machine. Actual type of the address depends on the adapter type. Typically either the public or external IP address.

* `created_at` - Date when entity was created. Date and time format is ISO 8601 and UTC.

* `external_id` - External entity ID on provider side.

* `external_region_id` - External regionId of resource.

* `external_zone_id` - External zoneId of resource.

* `organization_id` - ID of organization that entity belongs to.

* `owner` - Email of entity owner.

* `power_state` - Power state of machine.

* `updated_at` - Date when entity was last updated. Date and time format is ISO 8601 and UTC.

