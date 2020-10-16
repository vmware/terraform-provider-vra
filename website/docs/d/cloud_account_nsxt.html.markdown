---
layout: "vra"
page_title: "VMware vRealize Automation: vra_cloud_account_nsxt"
description: |-
    Provides a data lookup for vra_cloud_account_nsxt.
---

# Data Source: vra\_cloud\_account\_nsxt

Provides a VMware vRA vra_cloud_account_nsxt data source.

## Example Usages

**NSX-T cloud account data source by its id:**

This is an example of how to create a NSX-T cloud account resource and read it as a data source using its id.
NOTE: The NSX-T cloud account resource need not be created through terraform.
To create a NSX-T cloud account, follow the resource NSX-T cloud account documentation:

```hcl

resource "vra_cloud_account_nsxt" "this" {
  name        = "tf-nsx-t-account"
  description = "foobar"
  username    = var.username
  password    = var.password
  hostname    = var.hostname
  dc_id       = var.cloud_proxy != "" ? data.vra_data_collector.dc[0].id : ""

  accept_self_signed_cert = true

  tags {
    key   = "foo"
    value = "bar"
  }
}

data "vra_cloud_account_nsxt" "this" {
  id = "vra_cloud_account_nsxt.this.id"
}

```

**NSX-T cloud account data source by its name:**

This is an example of how to create a NSX-T cloud account resource and read it as a data source using its name.
NOTE: The NSX-T cloud account resource need not be created through terraform.
To create a NSX-T cloud account, follow the resource NSX-T cloud account documentation:

```hcl

resource "vra_cloud_account_nsxt" "this" {
  name        = "tf-nsx-t-account"
  description = "foobar"
  username    = var.username
  password    = var.password
  hostname    = var.hostname
  dc_id       = var.cloud_proxy != "" ? data.vra_data_collector.dc[0].id : ""

  accept_self_signed_cert = true

  tags {
    key   = "foo"
    value = "bar"
  }
}

data "vra_cloud_account_nsxt" "this" {
  name = "vra_cloud_account_nsxt.this.name"
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

* `links` - HATEOAS of the entity.

* `org_id` - The id of the organization this entity belongs to.

* `owner` - Email of the user that owns the entity.

* `tags` - A set of tag keys and optional values that were set on this resource.
example:[ { "key" : "vmware", "value": "provider" } ]
  * `key` - Tag’s key.
  * `value` - Tag’s value.

* `updated_at` - Date when the entity was last updated. The date is ISO 8601 and UTC.

* `username` - Username to authenticate with the cloud account.
