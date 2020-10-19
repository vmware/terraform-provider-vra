---
layout: "vra"
page_title: "VMware vRealize Automation: vra_network"
description: |-
  Provides a VMware vRA vra_network resource.
---

# Resource: vra_machine
## Example Usages

This is an example of how to create a network resource.

```hcl
resource "vra_network" "my_network" {
  name = "terraform_vra_network-%d"
  outbound_access = false
  tags {
	key = "stoyan"
    value = "genchev"
  }
  constraints {
	  mandatory = true
	  expression = "pci"
  }
}
```
A network resource supports the following resource:

## Argument Reference

* `id` - The id of the image profile instance.  Only one of 'filter', 'id', 'name' or 'region_id' must be specified.

* `name` - A human-friendly name used as an identifier in APIs that support this option.  Only one of 'filter', 'id', 'name' or 'region_id' must be specified.

## Attribute Reference

* `cidr` - IPv4 address range of the network in CIDR format.

* `constraints` - List of storage, network and extensibility constraints to be applied when provisioning through this project.

* `custom_properties` - Additional properties that may be used to extend the base resource.

* `deployment_id` - Deployment id that is associated with this resource.

* `description` - A human-friendly description.

* `external_id` - External entity Id on the provider side.

* `external_zone_id` - The external zoneId of the resource.

* `links` - HATEOAS of the entity

* `organization_id` - The id of the organization this entity belongs to.

* `outbound_access` - Flag to indicate if the network needs to have outbound access or not. Default is true. This field will be ignored if there is proper input for networkType customProperty

* `owner` - Email of the user that owns the entity.

* `project_id` - The id of the project this resource belongs to.

* `self_link` - Self link of this request.

* `tags` - A set of tag keys and optional values that were set on this resource.
           example:[ { "key" : "ownedBy", "value": "Rainpole" } ]

* `update_at` - Date when the entity was last updated. The date is ISO 8601 and UTC.
