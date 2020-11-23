---
layout: "vra"
page_title: "VMware vRealize Automation: fabric_network"
description: |-
  Provides a data lookup for vRA fabric networks.
---

# Data Source: fabric_network
## Example Usages
This is an example of how to lookup fabric networks.

**Fabric network data source by Id:**

```hcl
# Lookup vRA fabric network using its Id
data "fabric_network" "this" {
  id = var.fabric_network_id
}
```

**Fabric network by filter query:**

```hcl
# Lookup vRA fabric network using its name
data "fabric_network" "this" {
  filter = "name eq '${var.name}'"
}
```

A fabric network data source supports the following arguments:

## Argument Reference
* `filter` - Filter query string that is supported by vRA multi-cloud IaaS API. Only one of 'filter' or 'id' must be specified.

* `id` - The id of the vRA fabric network.  Only one of 'filter' or 'id' must be specified.

## Attribute Reference
* `cloud_account_ids` - Set of ids of the cloud accounts this entity belongs to.

* `created_at` - Date when the entity was created. The date is in ISO 6801 and UTC.

* `cidr` - Network CIDR to be used.

* `is_public` - Indicates whether the sub-network supports public IP assignment.

* `is_default` - Indicates whether this is the default subnet for the zone.

* `description` - State object representing a network on a external cloud provider.

* `external_id` - External entity Id on the provider side.

* `external_region_id` - The id of the region for which this network is defined.

* `links` - HATEOAS of the entity

* `name` - A human-friendly name used as an identifier in APIs that support this option.

* `organization_id` - The id of the organization this entity belongs to.

* `owner` - Email of the user that owns the entity.

* `tags` -  A set of tag keys and optional values that were set on this resource.
                       example:[ { "key" : "ownedBy", "value": "Rainpole" } ]

* `updated_at` - Date when the entity was last updated. The date is ISO 8601 and UTC.