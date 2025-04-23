---
page_title: "VMware Aria Automation: vra_network"
description: |-
  Provides a data lookup for vra_network.
---

# Data Source: vra_network

## Example Usages

This is an example of how to read a network data source.

**Network data source by id:**

```hcl
data "vra_network" "test-network" {
  id = var.network_id
}
```

**Network data source by name:**

```hcl
data "vra_network" "test-network" {
  name = var.network_name
}
```

**Network data source by filter:**

```hcl
data "vra_network" "test-network" {
  filter = "name eq '${var.network_name}'"
}
```

## Argument Reference

* `id` - (Optional) The id of the network instance. Only one of `id`, `name`, or `filter` must be specified.

* `name` - (Optional) The human-friendly name of the network instance. Only one of `id`, `name`, or `filter` must be specified.

* `filter` - (Optional) The search criteria to narrow down the network instance. Only one of `id`, `name`, or `filter` must be specified.

## Attribute Reference

* `cidr` - IPv4 address range of the network in CIDR format.

* `cloud_account_ids` - Set of ids of the cloud accounts this resource belongs to.

* `custom_properties` - Additional properties that may be used to extend the base resource.

* `deployment_id` - Deployment id that is associated with this resource.

* `description` - A human-friendly description.

* `external_id` - External entity Id on the provider side.

* `external_region_id` - The external regionId of the resource.

* `external_zone_id` - The external zoneId of the resource.

* `links` - Hypermedia as the Engine of Application State (HATEOAS) of the entity.

* `organization_id` - The id of the organization this entity belongs to.

* `owner` - Email of the user or display name of the group that owns the entity.

* `project_id` - The id of the project this resource belongs to.

* `tags` - A set of tag keys and optional values that were set on this resource. Example: `[ { "key" : "ownedBy", "value": "Rainpole" } ]`

* `update_at` - Date when the entity was last updated. The date is ISO 8601 and UTC.
