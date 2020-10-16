---
layout: "vra"
page_title: "VMware vRealize Automation: vra_cloud_account_vmc"
description: |-
    Provides a data lookup for vra_cloud_account_vmc.
---

# Data Source: vra\_cloud\_account\_vmc

Provides a VMware vRA vra_cloud_account_vmc data source.

## Example Usages

**VMC cloud account data source by its id:**

This is an example of how to create a vmc cloud account resource and read it as a data source using its id.
NOTE: The vmc cloud account resource need not be created through terraform.
To create a vmc cloud account, follow the resource vmc cloud account documentation:

```hcl

// Required for vRA Cloud, Optional for vRA on-prem
data "vra_data_collector" "this" {
  count = var.data_collector_name != "" ? 1 : 0
  name  = var.data_collector_name
}

data "vra_region_enumeration_vmc" "this" {
  api_token = var.api_token
  sddc_name = var.sddc_name

  vcenter_hostname = var.vcenter_hostname
  vcenter_password = var.vcenter_password
  vcenter_username = var.vcenter_username
  dc_id            = var.data_collector_name != "" ? data.vra_data_collector.this[0].id : "" // Required for vRA Cloud, Optional for vRA on-prem
}

resource "vra_cloud_account_vmc" "this" {
  name        = "tf-vra-cloud-account-vmc"
  description = "tf test vmc cloud account"
  api_token = var.api_token
  sddc_name = var.sddc_name

  vcenter_hostname = var.vcenter_hostname
  vcenter_password = var.vcenter_password
  vcenter_username = var.vcenter_username
  nsx_hostname     = var.nsx_hostname
  dc_id            = var.data_collector_name != "" ? data.vra_data_collector.this[0].id : "" // Required for vRA Cloud, Optional for vRA on-prem

  regions                 = data.vra_cloud_account_vmc.this.regions
  accept_self_signed_cert = true

  tags {
    key   = "foo"
    value = "bar"
  }
}

data "vra_cloud_account_vmc" "this" {
  id = "vra_cloud_account_vmc.this.id"
}

```

**vmc cloud account data source by its name:**

This is an example of how to create a vmc cloud account resource and read it as a data source using its name.
NOTE: The vmc cloud account resource need not be created through terraform.
To create a vmc cloud account, follow the resource vmc cloud account documentation:

```hcl

data "vra_data_collector" "dc" {
  count = var.cloud_proxy != "" ? 1 : 0
  name  = var.cloud_proxy
}

data "vra_region_enumeration_vmc" "this" {
  username                = var.username
  password                = var.password
  hostname                = var.hostname
  dcid                    = var.cloud_proxy != "" ? data.vra_data_collector.dc[0].id : "" // Required for vRA Cloud, Optional for vRA 8.X
  accept_self_signed_cert = true
}

resource "vra_cloud_account_vmc" "this" {
  name        = "tf-vra-cloud-account-vmc"
  description = "tf test vmc cloud account"

  api_token = var.api_token
  sddc_name = var.sddc_name

  vcenter_hostname = var.vcenter_hostname
  vcenter_password = var.vcenter_password
  vcenter_username = var.vcenter_username
  nsx_hostname     = var.nsx_hostname
  dc_id            = var.data_collector_name != "" ? data.vra_data_collector.this[0].id : "" // Required for vRA Cloud, Optional for vRA on-prem

  regions                 = data.vra_region_enumeration_vmc.this.regions
  accept_self_signed_cert = true

  tags {
    key   = "foo"
    value = "bar"
  }
}

data "vra_cloud_account_vmc" "this" {
  name = "vra_cloud_account_vmc.this.name"
}

```



## Argument Reference

The following arguments are supported for an vmc cloud account data source:

* `id` - (Optional) The id of this vmc cloud account.

* `name` - (Optional) The name of this vmc cloud account.

## Attribute Reference

* `created_at` - Date when the entity was created. The date is in ISO 6801 and UTC.

* `dc_id` - Identifier of a data collector vm deployed in the on premise infrastructure. Refer to the data-collector API to create or list data collector.

* `description` - A human-friendly description.

* `links` - HATEOAS of the entity.

* `nsx_hostname` - The IP address of the NSX Manager server in the specified SDDC / FQDN.

* `org_id` - The id of the organization this entity belongs to.

* `owner` - Email of the user that owns the entity.

* `regions` - A set of region names that are enabled for this account.

* `region-ids` - A set of region IDs that are enabled for this account.

* `sddc_name` - Identifier of the on-premise SDDC to be used by this cloud account. Note that NSX-V SDDCs are not supported.

* `tags` - A set of tag keys and optional values that were set on this resource.
example:[ { "key" : "vmware", "value": "provider" } ]
  * `key` - Tag’s key.
  * `value` - Tag’s value.
  
* `updated_at` - Date when the entity was last updated. The date is ISO 8601 and UTC.

* `vcenter_hostname` - The IP address or FQDN of the vCenter Server in the specified SDDC. The cloud proxy belongs on this vCenter.

* `vcenter_username` - vCenter user name for the specified SDDC. The specified user requires CloudAdmin credentials. The user does not require CloudGlobalAdmin credentials.