---
layout: "vra"
page_title: "VMware vRealize Automation: vra_cloud_account_azure"
description: |-
    Provides a data lookup for vra_cloud_account_azure.
---

# Data Source: vra\_cloud\_account\_azure

Provides a VMware vRA vra_cloud_account_azure data source.

## Example Usages

**Azure cloud account data source by its id:**

This is an example of how to read the cloud account data source using its id.

```hcl

data "vra_cloud_account_azure" "this" {
  id = var.vra_cloud_account_azure_id
}

```

**Azure cloud account data source by its name:**

This is an example of how to read the cloud account data source using its name.

```hcl

data "vra_cloud_account_azure" "this" {
  name = var.vra_cloud_account_azure_name
}

```



## Argument Reference

The following arguments are supported for an Azure cloud account data source:

* `id` - (Optional) The id of this Azure cloud account.

* `name` - (Optional) The name of this Azure cloud account.

## Attribute Reference

* `application_id` - Azure Client Application ID.

* `created_at` - Date when the entity was created. The date is in ISO 6801 and UTC.

* `description` - A human-friendly description.

* `links` - HATEOAS of the entity.

* `org_id` - The id of the organization this entity belongs to.

* `owner` - Email of the user that owns the entity.

* `regions` - A set of region names that are enabled for this account.

* `subscription_id` - Azure Subscription ID.

* `tags` - A set of tag keys and optional values that were set on this resource.
example: `[ { "key" : "vmware", "value": "provider" } ]`
  * `key` - Tag’s key.
  * `value` - Tag’s value.

* `tenant_id` - Azure Tenant ID.

* `updated_at` - Date when the entity was last updated. The date is ISO 8601 and UTC.
