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

    * `content` - (Optional) Calid cloud config data in json-escaped yaml syntax.

* `custom_properties` - (Optional) Additional properties that may be used to extend the base resource.

* `deployment_id` - (Optional) Describes machine within the scope of your organization and is not propagated to the cloud.

* `description` - (Optional) A human-friendly description.

* `disks` - (Optional) Specification for attaching/detaching disks to a machine.

    * `block_device_id` - (Required) ID of the existing block device.

    * `description` - (Optional) Human-friendly description.

    * `name` - (Optional) Human-friendly block-device name used as an identifier in APIs that support this option.

* `flavor` - (Required) Flavor of machine instance.

* `image` - (Optional) Type of image used for this machine.

* `image_disk_constraints` - (Optional) Constraints that are used to drive placement policies for the image disk. Constraint expressions are matched against tags on existing placement targets. example: `[{"mandatory" : "true", "expression": "environment:prod"}, {"mandatory" : "false", "expression": "pci"}]`. It is nested argument with the following properties.

    * `expression` - (Required) Constraint that is conveyed to the policy engine. An expression of the form "[!]tag-key[:[tag-value]]", used to indicate a constraint match on keys and values of tags.

    * `mandatory` - (Required) Indicates whether this constraint should be strictly enforced or not.

* `image_ref` - (Optional) Direct image reference used for this machine (name, path, location, uri, etc.). Valid if no image property is provided

* `name` - (Required) Human-friendly name used as an identifier in APIs that support this option.

* `nics` - (Optional) Set of network interface controller specifications for this machine. If not specified, then a default network connection will be created.

    * `addresses` - (Optional) List of IP addresses allocated or in use by this network interface.
                    example: `[ "10.1.2.190" ]`

    * `custom_properties` - (Optional) Additional properties that may be used to extend the base type.

    * `description` - (Optional) Human-friendly description.

    * `device_index` - (Optional) The device index of this network interface.

    * `name` - (Optional) Human-friendly name used as an identifier in APIs that support this option.

    * `network_id` - (Required) ID of the network instance that this network interface plugs into.

    * `security_group_ids` - (Optional) List of security group ids which this network interface will be assigned to.

* `tags` - (Optional) Set of tag keys and optional values that should be set on any resource that is produced from this specification. example: `[ { "key" : "ownedBy", "value": "Rainpole" } ]`. It is nested argument with the following properties.

    * `key` - (Required) Tag’s key.

    * `value` - (Required) Tag’s value.

## Attribute Reference

* `address` - Primary address allocated or in use by this machine. The actual type of the address depends on the adapter type. Typically it is either the public or the external IP address.

* `constraints` - Constraints used to drive placement policies for the virtual machine produced from the specification. Constraint expressions are matched against tags on existing placement targets.
Example: `[{"mandatory" : "true", "expression": "environment":"prod"}, {"mandatory" : "false", "expression": "pci"}]`

* `created_at` - Date when the entity was created. Date and time format is ISO 8601 and UTC.

* `attach_disks_before_boot` - By default, disks are attached after the machine has been built. FCDs cannot be attached to machine as a day 0 task.

* `disks_list` - List of all disks attached to a machine including boot disk, and additional block devices attached using the disks attribute.

    * `block_device_id` - ID of existing block device.

    * `description` - Human-friendly description.

    * `name` - Human-friendly block-device name used as an identifier in APIs that support this option.

    * `scsi_controller` - The id of the SCSI controller (_e.g_., `SCSI_Controller_0`.)

    * `unit_number` - The unit number of the SCSI controller (_e.g_., `2`.)

* `external_id` - External entity ID on the provider side.

* `external_region_id` - External regionId of the resource.

* `external_zone_id` - External zoneId of the resource.

* `links` - HATEOAS of the entity

* `organization_id` - ID of the organization this entity belongs to.

* `owner` - Email of entity owner.

* `power_state` - Power state of machine.

* `project_id` - ID of project that resource belongs to.

* `update_at` - Date when the entity was last updated. The date is ISO 8601 and UTC.
