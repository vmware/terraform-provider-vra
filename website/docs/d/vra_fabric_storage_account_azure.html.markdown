---
layout: "vra"
page_title: "VMware vRealize Automation: vra_fabric_storage_account_azure"
description: |-
  Provides a data lookup for fabric Azure storage account.
---

# Data Source: vra_fabric_storage_account_azure
## Example Usages
This is an example of how to lookup fabric Azure storage account.

**Fabric Azure storage account by Id:**

```hcl
# Lookup fabric Azure storage account using its Id
data "vra_fabric_storage_account_azure" "this" {
  id = var.fabric_storage_account_azure_id
}
```

**Fabric Azure storage by filter query:**

```hcl
# Lookup fabric Azure storage account using its name
data "vra_fabric_storage_account_azure" "this" {
  filter = "name eq '${var.name}'"
}
```

A fabric Azure storage account supports the following arguments:

## Argument Reference
* `filter` - Search criteria to narrow down the fabric Azure storage accounts. Only one of 'filter' or 'id' must be specified.

* `id` - The id of the fabric Azure storage account. Only one of 'filter' or 'id' must be specified.

## Attribute Reference
* `cloud_account_ids` - Set of ids of the cloud accounts this entity belongs to.

* `created_at` - Date when the entity was created. The date is in ISO 6801 and UTC.

* `description` - A human-friendly description of the fabric Azure storage account.

* `external_id` - External entity Id on the provider side.

* `external_region_id` - The id of the region for which this entity is defined.

* `links` - HATEOAS of the entity

* `name` - A human-friendly name used as an identifier in APIs that support this option.

* `org_id` - The id of the organization this entity belongs to.

* `owner` - Email of the user that owns the entity.

* `tags` -  A set of tag keys and optional values that were set on this resource.
                       example: `[ { "key" : "ownedBy", "value": "Rainpole" } ]`

* `type` -  Indicates the performance tier for the storage type. Premium disks are SSD backed and Standard disks are HDD backed. example: `Standard_LRS / Premium_LRS`

* `updated_at` - Date when the entity was last updated. The date is ISO 8601 and UTC.
