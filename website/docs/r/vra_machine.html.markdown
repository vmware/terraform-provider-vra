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

* `block_device_id` - The id of the block device.

* `description` - A human-friendly description.

* `filter` - Filter query string that is supported by vRA multi-cloud IaaS API. Example: regionId eq '<regionId>' and cloudAccountId eq '<cloudAccountId>'. Only one of 'filter', 'id', 'name' or 'region_id' must be specified.

* `id` - The id of the image profile instance.  Only one of 'filter', 'id', 'name' or 'region_id' must be specified.

## Attribute Reference

* `address` - Primary address allocated or in use by this machine. The actual type of the address depends on the adapter type. Typically it is either the public or the external IP address.

* `cloud_account_ids` - Set of ids of the cloud accounts this resource belongs to.

* `created_at` - Date when the entity was created. The date is in ISO 6801 and UTC.

* `custom_properties` - Additional properties that may be used to extend the base resource.

* `deployment_id` - Deployment id that is associated with this resource.

* `external_id` - External entity Id on the provider side.

* `external_region_id` - The external regionId of the resource.

* `external_zone_id` - The external zoneId of the resource.

* `links` - HATEOAS of the entity

* `name` - A human-friendly name used as an identifier in APIs that support this option.  Only one of 'filter', 'id', 'name' or 'region_id' must be specified.

* `org_id` - The id of the organization this entity belongs to.

* `owner` - Email of the user that owns the entity.

* `power_state` - Power state of machine.

* `project_id` - The id of the project this resource belongs to.

* `tags` - A set of tag keys and optional values that were set on this resource.
           example:[ { "key" : "ownedBy", "value": "Rainpole" } ]

* `update_at` - Date when the entity was last updated. The date is ISO 8601 and UTC.
