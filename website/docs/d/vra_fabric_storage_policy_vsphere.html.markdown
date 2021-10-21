---
layout: "vra"
page_title: "VMware vRealize Automation: vra_fabric_storage_policy_vsphere"
description: |-
  Provides a data lookup for fabric vSphere storage policy.
---

# Data Source: vra_fabric_storage_policy_vsphere
## Example Usages
This is an example of how to lookup fabric vSphere storage policies.

**Fabric vSphere storage policy by Id:**

```hcl
# Lookup fabric vSphere storage account using its Id
data "vra_fabric_storage_policy_vsphere" "this" {
  id = var.fabric_storage_policy_vsphere_id
}
```

**Fabric vSphere storage policy by filter query:**

```hcl
# Lookup fabric vSphere storage policy using its name
data "vra_fabric_storage_policy_vsphere" "this" {
  filter = "name eq '${var.name}'"
}
```

A fabric vSphere storage policy supports the following arguments:

## Argument Reference
* `filter` - Search criteria to narrow down the fabric vSphere storage policy. Only one of 'filter' or 'id' must be specified.

* `id` - The id of the fabric vSphere storage policy. Only one of 'filter' or 'id' must be specified.

## Attribute Reference
* `cloud_account_ids` - Set of ids of the cloud accounts this entity belongs to.

* `created_at` - Date when the entity was created. The date is in ISO 6801 and UTC.

* `description` - A human-friendly description of the fabric vSphere storage account.

* `external_id` - External entity Id on the provider side.

* `external_region_id` - The id of the region for which this entity is defined.

* `links` - HATEOAS of the entity

* `name` - A human-friendly name used as an identifier in APIs that support this option.  Only one of 'filter', 'id', 'name' or 'region_id' must be specified.

* `org_id` - The id of the organization this entity belongs to.

* `updated_at` - Date when the entity was last updated. The date is ISO 8601 and UTC.