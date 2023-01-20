---layout: "vra"
page_title: "VMware vRealize Automation: vra_network"
description: |-
  Provides a data lookup for vra_network.
---

# Data Source: vra_network
## Example Usages

This is an example of how to read a network data source.

```hcl

data "vra_network" "test-network" {
  name = var.network_name
}

```

## Argument Reference

* `id` - (Optional) The id of the image profile instance.

* `name` - (Optional) A human-friendly name used as an identifier in APIs that support this option.

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
           example: `[ { "key" : "ownedBy", "value": "Rainpole" } ]`

* `update_at` - Date when the entity was last updated. The date is ISO 8601 and UTC.
