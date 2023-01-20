---
layout: "vra"
page_title: "VMware vRealize Automation: vra_cloud_account_nsxt"
description: |-
    Creates a vra_cloud_account_nsxt resource.
---

# Resource: vra\_cloud\_account\_nsxt

Creates a VMware vRealize Automation NSX-T cloud account resource.

## Example Usages

The following example shows how to create an NSX-T cloud account resource.

```hcl
resource "vra_cloud_account_nsxt" "this" {
  name        = "tf-nsx-t-account"
  description = "foobar"
  username    = var.username
  password    = var.password
  hostname    = var.hostname
  dc_id       = var.vra_data_collector_id

  accept_self_signed_cert = true

  tags {
    key   = "foo"
    value = "bar"
  }
}
```

## Argument Reference

Create your NSX-T cloud account resource with the following arguments:

* `accept_self_signed_cert` - (Optional) Accept self-signed certificate when connecting to the cloud account.

* `dc_id` - (Optional) Identifier of a data collector VM deployed in the on premise infrastructure.

* `description` - (Optional) Human-friendly description.

* `hostname` - (Required) Host name for NSX-T cloud account.

* `name` - (Optional) Name of NSX-T cloud account.

* `password` - (Required) Password used to authenticate to the cloud Account.

* `tags` - (Optional) Set of tag keys and values to apply to the cloud account.
Example: `[ { "key" : "vmware", "value": "provider" } ]`

* `username` - (Required) Username used to authenticate to the cloud account.

## Attribute Reference

* `associated_cloud_account_ids` - Cloud accounts associated with the cloud account.

* `created_at` - Date when entity was created. Date and time format is ISO 8601 and UTC.

* `id` - ID of NSX-T cloud account.

* `links` - HATEOAS of entity.

* `org_id` - ID of organization that entity belongs to.

* `owner` - Email of entity owner.

* `updated_at` - Date when entity was last updated. Date and time format is ISO 8601 and UTC.


## Import

To import the NSX-T cloud account, use the ID as in the following example:

`$ terraform import vra_cloud_account_nsxt.new_gcp 05956583-6488-4e7d-84c9-92a7b7219a15`
