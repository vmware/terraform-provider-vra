---
page_title: "VMware Aria Automation: vra_cloud_account_nsxt"
description: |-
    Provides a data lookup for vra_cloud_account_nsxt.
---

# Data Source: vra_cloud_account_nsxt

Provides a vra_cloud_account_nsxt data source.

## Example Usages

**NSX-T cloud account data source by its id:**

This is an example of how to read the cloud account data source using its id.

```hcl

data "vra_cloud_account_nsxt" "this" {
  id = var.vra_cloud_account_nsxt_id
}
```

**NSX-T cloud account data source by its name:**

This is an example of how to read the cloud account data source using its name.

```hcl

data "vra_cloud_account_nsxt" "this" {
  name = var.vra_cloud_account_nsxt_name
}
```

## Argument Reference

The following arguments are supported for an NSX-T cloud account data source:

* `id` - (Optional) The id of this NSX-T cloud account.

* `name` - (Optional) The name of this NSX-T cloud account.

## Attribute Reference

* `associated_cloud_account_ids` - Cloud accounts associated with this cloud account.

* `created_at` - Date when the entity was created. The date is in ISO 6801 and UTC.

* `dc_id` - Identifier of a data collector vm deployed in the on premise infrastructure.

* `description` - A human-friendly description.

* `hostname` - Host name for the NSX-T cloud account.

* `links` - Hypermedia as the Engine of Application State (HATEOAS) of the entity.

* `org_id` - The id of the organization this entity belongs to.

* `owner` - Email of the user that owns the entity.

* `tags` - A set of tag keys and optional values that were set on this resource. Example: `[ { "key" : "vmware", "value": "provider" } ]`

  * `key` - Tag’s key.

  * `value` - Tag’s value.

* `updated_at` - Date when the entity was last updated. The date is ISO 8601 and UTC.

* `username` - Username to authenticate with the cloud account.
