---
layout: "vra"
page_title: "VMware vRealize Automation: vra_cloud_account_nsxv"
description: |-
    Provides a data lookup for vra_cloud_account_nsxv.
---

# Data Source: vra\_cloud\_account\_nsxv

Provides a VMware vRA vra_cloud_account_nsxv data source.

## Example Usages

**NSX-V cloud account data source by its id:**

This is an example of how to read the cloud account data source using its id.

```hcl

data "vra_cloud_account_nsxv" "this" {
  id = var.vra_cloud_account_nsxv_id
}

```

**NSX-V cloud account data source by its name:**

This is an example of how to read the cloud account data source using its name.

```hcl

data "vra_cloud_account_nsxv" "this" {
  name = var.vra_cloud_account_nsxv_name
}

```



## Argument Reference

The following arguments are supported for an NSX-V cloud account data source:

* `id` - (Optional) The id of this NSX-V cloud account.

* `name` - (Optional) The name of this NSX-V cloud account.

## Attribute Reference

* `associated_cloud_account_ids` - Cloud accounts associated with this cloud account.

* `created_at` - Date when the entity was created. The date is in ISO 6801 and UTC.

* `dc_id` - Identifier of a data collector vm deployed in the on premise infrastructure.

* `description` - A human-friendly description.

* `hostname` - Host name for the NSX-V cloud account.

* `links` - HATEOAS of the entity.

* `org_id` - The id of the organization this entity belongs to.

* `owner` - Email of the user that owns the entity.

* `tags` - A set of tag keys and optional values that were set on this resource.
example: `[ { "key" : "vmware", "value": "provider" } ]`
  * `key` - Tag’s key.
  * `value` - Tag’s value.

* `updated_at` - Date when the entity was last updated. The date is ISO 8601 and UTC.

* `username` - Username to authenticate with the cloud account.
