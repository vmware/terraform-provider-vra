---
layout: "vra"
page_title: "VMware vRealize Automation: vra_fabric_network"
description: |-
  Provides a data lookup for vRA fabric networks.
---

# Data Source: vra_fabric_network
## Example Usages
This is an example of how to lookup fabric networks.

**Fabric network by filter query:**

```hcl
# Lookup vRA fabric network using its name
data "vra_fabric_network" "this" {
  filter = "name eq '${var.name}'"
}

# Lookup vRA fabric network using its name and regionId
data "vra_fabric_network" "this" {
  filter = "name eq '${var.name}' and externalRegionId eq '${var.external_region_id}'"
}
```

A fabric network data source supports the following arguments:

## Argument Reference
* `filter` - (Required) Filter query string that is supported by vRA multi-cloud IaaS API.

## Attribute Reference
* `cloud_account_ids` - Set of ids of the cloud accounts this entity belongs to.

* `created_at` - Date when the entity was created. The date is in ISO 6801 and UTC.

* `cidr` - Network CIDR to be used.

* `id` - ID of the vRA fabric network.

* `is_public` - Indicates whether the sub-network supports public IP assignment.

* `is_default` - Indicates whether this is the default subnet for the zone.

* `description` - State object representing a network on a external cloud provider.

* `external_id` - External entity Id on the provider side.

* `external_region_id` - The id of the region for which this network is defined.

* `links` - HATEOAS of the entity

* `name` - Name of the fabric network.

* `organization_id` - ID of organization that entity belongs to.

* `owner` - Email of the user that owns the entity.

* `tags` -  Set of tag keys and values to apply to the resource.
            Example: `[ { "key" : "vmware", "value": "provider" } ]`

* `updated_at` - Date when the entity was last updated. The date is ISO 8601 and UTC.

* `custom_properties` - Additional properties that may be used to extend the base resource.
