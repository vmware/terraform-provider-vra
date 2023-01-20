---
layout: "vra"
page_title: "VMware vRealize Automation: vra_network"
description: |-
  Provides a VMware vRA vra_network resource.
---

# Resource: vra_network
## Example Usages

This is an example of how to create a network resource.

```hcl
resource "vra_network" "my_network" {
  name = "terraform_vra_network-%d"
  outbound_access = false
  tags {
	key = "foo"
    value = "bar"
  }
  constraints {
	  mandatory = true
	  expression = "pci"
  }
}
```
A network resource supports the following resource:

## Argument Reference

* `custom_properties` - (Optional) Additional properties that may be used to extend the base resource.

* `deployment_id` - (Optional) Deployment id that is associated with this resource.

* `description` - (Optional) A human-friendly description.

* `name` - (Required) A human-friendly name used as an identifier in APIs that support this option.

* `outbound_access` - (Optional) Flag to indicate if the network needs to have outbound access or not. Default is true. This field will be ignored if there is proper input for networkType customProperty.

* `project_id` - (Required) The id of the project this resource belongs to.

## Attribute Reference

* `cidr` - IPv4 address range of the network in CIDR format.

* `constraints` - List of storage, network and extensibility constraints to be applied when provisioning through this project.

* `external_id` - External entity Id on the provider side.

* `external_zone_id` - The external zoneId of the resource.

* `links` - HATEOAS of the entity

* `organization_id` - The id of the organization this entity belongs to.

* `owner` - Email of the user that owns the entity.

* `self_link` - Self link of this request.

* `tags` - A set of tag keys and optional values that were set on this resource.
           example: `[ { "key" : "ownedBy", "value": "Rainpole" } ]`

* `update_at` - Date when the entity was last updated. The date is ISO 8601 and UTC.
