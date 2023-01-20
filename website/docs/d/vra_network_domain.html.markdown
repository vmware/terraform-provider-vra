---
layout: "vra"
page_title: "VMware vRealize Automation: vra_network_domain"
description: |-
  Provides a data lookup for Network domain objects.
---

# Data Source: vra_network_domain
## Example Usages
This is an example of how to lookup Network domain objects.

**Network Domain by filter query:**

```hcl
# Lookup network domain using its name
data "vra_network_domain" "this" {
  filter = "name eq '${var.name}'"
}
```

A network domain object supports the following arguments:

## Argument Reference
* `filter` - (Required) Search criteria to narrow down the network domain objects.

## Attribute Reference
* `cloud_account_ids` - Set of ids of the cloud accounts this entity belongs to.

* `created_at` - Date when the entity was created. The date is in ISO 6801 and UTC.

* `description` - A human-friendly description of the fabric vSphere storage account.

* `external_id` - External entity Id on the provider side.

* `external_region_id` - The id of the region for which this entity is defined.

* `id` - ID of the fabric network domain object.

* `links` - HATEOAS of the entity

* `name` - Name of the network domain.

* `org_id` - ID of organization that entity belongs to.

* `owner` - Email of the user that owns the entity.

* `tags` -  Set of tag keys and values to apply to the resource.
            Example: `[ { "key" : "vmware", "value": "provider" } ]`

* `updated_at` - Date when the entity was last updated. The date is ISO 8601 and UTC.
