---
layout: "vra"
page_title: "VMware vRealize Automation: network_domain"
description: |-
  Provides a data lookup for Network domain objects.
---

# Data Source: network_domain
## Example Usages
This is an example of how to lookup Network domain objects.

**Network Domain by Id:**

```hcl
# Lookup network domain using its Id
data "network_domain" "this" {
  id = var.network_domain_id
}
```

**Network Domain by filter query:**

```hcl
# Lookup network domain using its name
data "network_domain" "this" {
  filter = "name eq '${var.name}'"
}
```

A network domain object supports the following arguments:

## Argument Reference
* `filter` - Search criteria to narrow down the network domain objects. Only one of 'filter' or 'id' must be specified.

* `id` - The id of the fabric network domain objects. Only one of 'filter' or 'id' must be specified.

## Attribute Reference
* `cloud_account_ids` - Set of ids of the cloud accounts this entity belongs to.

* `created_at` - Date when the entity was created. The date is in ISO 6801 and UTC.

* `description` - A human-friendly description of the fabric vSphere storage account.

* `external_id` - External entity Id on the provider side.

* `external_region_id` - The id of the region for which this entity is defined.

* `links` - HATEOAS of the entity

* `name` - A human-friendly name used as an identifier in APIs that support this option. 

* `org_id` - The id of the organization this entity belongs to.

* `owner` - Email of the user that owns the entity.

* `tags` -  A set of tag keys and optional values that were set on this resource.
                       example:[ { "key" : "ownedBy", "value": "Rainpole" } ]
                       
* `updated_at` - Date when the entity was last updated. The date is ISO 8601 and UTC.