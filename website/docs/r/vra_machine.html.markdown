---
layout: "vra"
page_title: "VMware vRealize Automation: vra_machine"
description: |-
  Creates a vra_machine resource.
---

# Resource: vra_machine

Creates a VMware vRealize Automation machine resource.

## Example Usages

The following example shows how to create a machine resource.

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

Create your machine resource with the following arguments:

* `boot_config` - (Optional)  Machine boot config that will be passed to the instance. Used to perform common automated configuration tasks and even run scripts after instance starts.
    
    * `content` - Valid cloud config data in JSON-escaped YAML syntax.
    
* `custom_properties` - (Optional) Additional properties that may be used to extend the base resource.

* `deployment_id` - (Optional) Describes machine within the scope of your organization and is not propagated to the cloud.

* `description` - (Optional) A human-friendly description.

* `disks` - (Optional) Specification for attaching/detaching disks to a machine.
    
    * `block_device_id` - ID of existing block device.
    
    * `description` - Human-friendly description.
    
    * `name` - Human-friendly block-device name used as an identifier in APIs that support this option.
    
* `flavor` - (Required) Flavor of machine instance.

* `image` - (Optional) Type of image used for this machine.

* `image_ref` - (Optional) Direct image reference used for this machine (name, path, location, uri, etc.). Valid if no image property is provided

* `name` - (Required) Human-friendly name used as an identifier in APIs that support this option.

## Attribute Reference

* `address` - Primary address allocated or in use by this machine. The actual type of the address depends on the adapter type. Typically it is either the public or the external IP address.

* `constraints` - Constraints used to drive placement policies for the virtual machine produced from the specification. Constraint expressions are matched against tags on existing placement targets.  
Example:[{"mandatory" : "true", "expression": "environment":"prod"}, {"mandatory" : "false", "expression": "pci"}]

* `created_at` - Date when the entity was created. Date and time format is ISO 8601 and UTC.

* `disks_list` - List of all disks attached to a machine including boot disk, and additional block devices attached using the disks attribute.
    
    * `block_device_id` - ID of existing block device.
    
    * `description` - Human-friendly description.
    
    * `name` - Human-friendly block-device name used as an identifier in APIs that support this option.
    
* `external_id` - External entity ID on the provider side.

* `external_region_id` - External regionId of the resource.

* `external_zone_id` - External zoneId of the resource.

* `image_disk_constraints` - Constraints used to drive placement policies for the image disk. Constraint expressions are matched against tags on existing placement targets.  
Example:[{"mandatory" : "true", "expression": "environment:prod"}, {"mandatory" : "false", "expression": "pci"}]

* `links` - HATEOAS of entity

* `nics` - Set of network interface controller specifications for this machine. If left unspecified, a default network connection is created.
    
    * `addresses` - List of IP addresses allocated or in use by this network interface.
                    example:[ "10.1.2.190" ]
    
    * `custom_properties` - Additional properties that may be used to extend the base type.
    
    * `description` - Human-friendly description.

    * `device_index` - Device index of this network interface.
    
    * `name` - Human-friendly name used as an identifier in APIs that support this option.
    
    * `network_id` - ID of network instance that network interface plugs into.

    * `security_group_ids` - List of security group IDs that network interface will be assigned to.
    
* `organization_id` - ID of organization that entity belongs to.

* `owner` - Email of entity owner.

* `power_state` - Power state of machine.

* `project_id` - ID of project that resource belongs to.

* `tags` - Set of tag keys and values set on the resource.
           Example:[ { "key" : "ownedBy", "value": "Rainpole" } ]

* `update_at` - Date when the entity was last updated. Date and time format is ISO 8601 and UTC.
