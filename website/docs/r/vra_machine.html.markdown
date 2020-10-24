---
layout: "vra"
page_title: "VMware vRealize Automation: vra_machine"
description: |-
  Provides a VMware vRA vra_machine resource.
---

# Resource: vra_machine
## Example Usages

This is an example on how to create a machine resource.

```hcl
resource "vra_machine" "this" {
  name        = "tf-machine"
  description = "terrafrom test machine"
  project_id  = data.vra_project.this.id
  image       = "ubuntu2"
  flavor      = "medium"

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
      ssh-rsa your-ssh-rsa:
    - sudo sed -e 's/.*PasswordAuthentication yes.*/PasswordAuthentication no/' -i /etc/ssh/sshd_config
    - sudo service sshd restart
EOF
  }

  nics {
    network_id = data.vra_network.this.id
  }

  constraints {
    mandatory  = true
    expression = "AWS"
  }

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
A machine resource supports the following resource:

## Argument Reference
* `boot_config` - (Optional)  Machine boot config that will be passed to the instance that can be used to perform common automated configuration tasks and even run scripts after the instance starts.
    
    * `content` - A valid cloud config data in json-escaped yaml syntax.
    
* `custom_properties` - (Optional) Additional properties that may be used to extend the base resource.

* `deployment_id` - (Optional) Describes machine within the scope of your organization and is not propagated to the cloud.

* `description` - (Optional) A human-friendly description.

* `disks` - (Optional) Specification for attaching/detaching disks to a machine.
    
    * `block_device_id` - The id of the existing block device.
    
    * `description` - A human-friendly description.
    
    * `name` - A human-friendly block-device name used as an identifier in APIs that support this option.
    
* `flavor` - (Required) Flavor of machine instance.

* `image` - (Optional) Type of image used for this machine.

* `image_ref` - (Optional) Direct image reference used for this machine (name, path, location, uri, etc.). Valid if no image property is provided

* `name` - (Required) A human-friendly name used as an identifier in APIs that support this option.

## Attribute Reference

* `address` - Primary address allocated or in use by this machine. The actual type of the address depends on the adapter type. Typically it is either the public or the external IP address.

* `constraints` - Constraints that are used to drive placement policies for the virtual machine that is produced from this specification. Constraint expressions are matched against tags on existing placement targets.
                  example:[{"mandatory" : "true", "expression": "environment":"prod"}, {"mandatory" : "false", "expression": "pci"}]

* `created_at` - Date when the entity was created. The date is in ISO 6801 and UTC.

* `disks_list` - List of all disks attached to a machine including boot disk, and additional block devices attached using the disks attribute.
    
    * `block_device_id` - The id of the existing block device.
    
    * `description` - A human-friendly description.
    
    * `name` - A human-friendly block-device name used as an identifier in APIs that support this option.
    
* `external_id` - External entity Id on the provider side.

* `external_region_id` - The external regionId of the resource.

* `external_zone_id` - The external zoneId of the resource.

* `image_disk_constraints` - Constraints that are used to drive placement policies for the image disk. Constraint expressions are matched against tags on existing placement targets.
                             example:[{"mandatory" : "true", "expression": "environment:prod"}, {"mandatory" : "false", "expression": "pci"}]

* `links` - HATEOAS of the entity

* `nics` - description:A set of network interface controller specifications for this machine. If not specified, then a default network connection will be created.
    
    * `addresses` - A list of IP addresses allocated or in use by this network interface.
                    example:[ "10.1.2.190" ]
    
    * `custom_properties` - Additional properties that may be used to extend the base type.
    
    * `description` - A human-friendly description.

    * `device_index` - The device index of this network interface.
    
    * `name` - A human-friendly name used as an identifier in APIs that support this option.
    
    * `network_id` - Id of the network instance that this network interface plugs into.

    * `security_group_ids` - A list of security group ids which this network interface will be assigned to.
    
* `organization_id` - The id of the organization this entity belongs to.

* `owner` - Email of the user that owns the entity.

* `power_state` - Power state of machine.

* `project_id` - The id of the project this resource belongs to.

* `tags` - A set of tag keys and optional values that were set on this resource.
           example:[ { "key" : "ownedBy", "value": "Rainpole" } ]

* `update_at` - Date when the entity was last updated. The date is ISO 8601 and UTC.
