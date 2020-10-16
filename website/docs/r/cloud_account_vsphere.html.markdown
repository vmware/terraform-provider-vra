---
layout: "vra"
page_title: "VMware vRealize Automation: vra_cloud_account_vsphere"
description: |-
    Provides a VMware vRA vra_cloud_account_vsphere resource.
---

# Data Source: vra\_cloud\_account\_vsphere

Provides a VMware vRA vra_cloud_account_vsphere resource.

## Example Usages

**Create vSphere cloud account:**

This is an example of how to create a vSphere cloud account resource.

```hcl

data "vra_data_collector" "dc" {
  count = var.cloud_proxy != "" ? 1 : 0
  name  = var.cloud_proxy
}

data "vra_region_enumeration_vsphere" "this" {
  username                = var.username
  password                = var.password
  hostname                = var.hostname
  dcid                    = var.cloud_proxy != "" ? data.vra_data_collector.dc[0].id : "" // Required for vRA Cloud, Optional for vRA 8.X
  accept_self_signed_cert = true
}

data "vra_cloud_account_nsxt" "this" {
  name        = var.vra_cloud_account_nsxt_name
}

resource "vra_cloud_account_vsphere" "this" {
  name        = "tf-vSphere-account"
  description = "foobar"
  username    = var.username
  password    = var.password
  hostname    = var.hostname
  dcid        = var.cloud_proxy != "" ? data.vra_data_collector.dc[0].id : "" // Required for vRA Cloud, Optional for vRA 8.X

  regions                      = data.vra_region_enumeration_vsphere.this.regions
  associated_cloud_account_ids = [data.vra_cloud_account_nsxt.this.id]

  accept_self_signed_cert = true

  tags {
    key   = "foo"
    value = "bar"
  }
}

```



## Argument Reference

The following arguments are supported for an vSphere cloud account resource:

* `accept_self_signed_cert` - (Optional) Accept self signed certificate when connecting to the cloud account.

* `dc_id` - (Optional) Identifier of a data collector vm deployed in the on premise infrastructure.

* `description` - (Optional) A human-friendly description.

* `hostname` - (Required) The IP address or FQDN of the vCenter Server. The cloud proxy belongs on this vCenter.

* `name` - (Optional) The name of this vSphere cloud account.

* `password` - (Required) Password for the user used to authenticate with the cloud Account.

* `regions` - (Optional) A set of region names that are enabled for this account.

* `tags` - (Optional) A set of tag keys and optional values that to set on this resource.
example:[ { "key" : "vmware", "value": "provider" } ]
  * `key` - Tag’s key.
  * `value` - Tag’s value.

* `username` - (Required) The vSphere username to authenticate the vSphere account.

## Attribute Reference

* `associated_cloud_account_ids` - Cloud accounts associated with this cloud account.

* `created_at` - Date when the entity was created. The date is in ISO 6801 and UTC.

* `id` - (Optional) The id of this vSphere cloud account.

* `links` - HATEOAS of the entity.

* `org_id` - The id of the organization this entity belongs to.

* `owner` - Email of the user that owns the entity.

* `region-ids` - A set of region IDs that are enabled for this account.

* `updated_at` - Date when the entity was last updated. The date is ISO 8601 and UTC.


## Import

VMC cloud account can be imported using the id, e.g.

`$ terraform import vra_cloud_account_vsphere.new_vsphere 05956583-6488-4e7d-84c9-92a7b7219a15`
