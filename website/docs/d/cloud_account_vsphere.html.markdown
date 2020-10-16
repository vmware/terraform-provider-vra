---
layout: "vra"
page_title: "VMware vRealize Automation: vra_cloud_account_vsphere"
description: |-
    Provides a data lookup for vra_cloud_account_vsphere.
---

# Data Source: vra\_cloud\_account\_vsphere

Provides a VMware vRA vra_cloud_account_vsphere data source.

## Example Usages

**vSphere cloud account data source by its id:**

This is an example of how to create a vSphere cloud account resource and read it as a data source using its id.
NOTE: The vSphere cloud account resource need not be created through terraform.
To create a vSphere cloud account, follow the resource vSphere cloud account documentation:

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

data "vra_cloud_account_vsphere" "this" {
  id = "vra_cloud_account_vsphere.this.id"
}

```

**vSphere cloud account data source by its name:**

This is an example of how to create a vSphere cloud account resource and read it as a data source using its name.
NOTE: The vSphere cloud account resource need not be created through terraform.
To create a vSphere cloud account, follow the resource vSphere cloud account documentation:

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
  name        = "tf-vsphere-account"
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

data "vra_cloud_account_vsphere" "this" {
  name = "vra_cloud_account_vsphere.this.name"
}

```



## Argument Reference

The following arguments are supported for an vSphere cloud account data source:

* `id` - (Optional) The id of this vSphere cloud account.

* `name` - (Optional) The name of this vSphere cloud account.

## Attribute Reference

* `associated_cloud_account_ids` - Cloud accounts associated with this cloud account.

* `created_at` - Date when the entity was created. The date is in ISO 6801 and UTC.

* `custom_properties` - Additional properties that may be used to extend the base info.

* `dc_id` - Identifier of a data collector vm deployed in the on premise infrastructure.

* `description` - A human-friendly description.

* `enabled_region_ids` - A set of region IDs that are enabled for this account.

* `hostname` - The IP address or FQDN of the vCenter Server. The cloud proxy belongs on this vCenter.

* `org_id` - The id of the organization this entity belongs to.

* `owner` - Email of the user that owns the entity.

* `tags` - A set of tag keys and optional values that were set on this resource.
example:[ { "key" : "vmware", "value": "provider" } ]
  * `key` - Tag’s key.
  * `value` - Tag’s value.

* `updated_at` - Date when the entity was last updated. The date is ISO 8601 and UTC.

* `username` - The vSphere username to authenticate the vsphere account.

