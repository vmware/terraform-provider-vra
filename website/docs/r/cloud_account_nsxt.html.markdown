---
layout: "vra"
page_title: "VMware vRealize Automation: vra_cloud_account_nsxt"
description: |-
    Provides a VMware vRA vra_cloud_account_nsxt resource.
---

# Data Source: vra\_cloud\_account\_nsxt

Provides a VMware vRA vra_cloud_account_nsxt resource.

## Example Usages

**Create NSX-T cloud account:**

This is an example of how to create a NSX-T cloud account resource.

```hcl

data "vra_data_collector" "dc" {
  count = var.cloud_proxy != "" ? 1 : 0
  name  = var.cloud_proxy
}

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

```



## Argument Reference

The following arguments are supported for an NSX-T cloud account resource:

* `accept_self_signed_cert` - (Optional) Accept self signed certificate when connecting to the cloud account.

* `dc_id` - (Optional) Identifier of a data collector vm deployed in the on premise infrastructure.

* `description` - (Optional) A human-friendly description.

* `hostname` - (Required) Host name for the NSX-T cloud account.

* `name` - (Optional) The name of this NSX-T cloud account.

* `password` - (Required)  Password for the user used to authenticate with the cloud Account.

* `tags` - (Optional) A set of tag keys and optional values that to set on this resource.
example:[ { "key" : "vmware", "value": "provider" } ]
  * `key` - Tag’s key.
  * `value` - Tag’s value.

* `username` - (Required) Username to authenticate with the cloud account.

## Attribute Reference

* `associated_cloud_account_ids` - Cloud accounts associated with this cloud account.

* `created_at` - Date when the entity was created. The date is in ISO 6801 and UTC.

* `id` - The id of this NSX-T cloud account.

* `links` - HATEOAS of the entity.

* `org_id` - The id of the organization this entity belongs to.

* `owner` - Email of the user that owns the entity.

* `updated_at` - Date when the entity was last updated. The date is ISO 8601 and UTC.


## Import

NSX-T cloud account can be imported using the id, e.g.

`$ terraform import vra_cloud_account_nsxt.new_gcp 05956583-6488-4e7d-84c9-92a7b7219a15`
