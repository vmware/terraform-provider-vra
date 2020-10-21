---
layout: "vra"
page_title: "VMware vRealize Automation: vra_cloud_account_vmc"
description: |-
    Provides a VMware vRA vra_cloud_account_vmc resource.
---

# Resource: vra\_cloud\_account\_vmc

Provides a VMware vRA vra_cloud_account_vmc resource.

## Example Usages

**Create VMC cloud account:**

This is an example of how to create a VMC cloud account resource.

```hcl

resource "vra_cloud_account_vmc" "this" {
  name        = "tf-vra-cloud-account-vmc"
  description = "tf test vmc cloud account"
  api_token = var.api_token
  sddc_name = var.sddc_name

  vcenter_hostname = var.vcenter_hostname
  vcenter_password = var.vcenter_password
  vcenter_username = var.vcenter_username
  nsx_hostname     = var.nsx_hostname
  dc_id            = var.data_collector_id  // Required for vRA Cloud, Optional for vRA on-prem
  regions                 = var.regions

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

* `api_token` - (Required) VMC API access key.

* `dc_id` - (Optional) Identifier of a data collector vm deployed in the on premise infrastructure. Refer to the data-collector API to create or list data collector.

* `description` - (Optional) A human-friendly description.

* `name` - (Optional) A human-friendly name used as an identifier in APIs that support this option.

* `nsx_hostname` - (Required)  The IP address of the NSX Manager server in the specified SDDC / FQDN.

* `regions` - (Optional) A set of region names that are enabled for this account.

* `sddc_name` - (Required) Identifier of the on-premise SDDC to be used by this cloud account. Note that NSX-V SDDCs are not supported.

* `tags` - (Optional) A set of tag keys and optional values that to set on this resource.
example:[ { "key" : "vmware", "value": "provider" } ]
  * `key` - Tag’s key.
  * `value` - Tag’s value.

* `vcenter_hostname` - (Required) The IP address or FQDN of the vCenter Server in the specified SDDC. The cloud proxy belongs on this vCenter.
  
* `vcenter_password` - (Required) Password for the user used to authenticate with the cloud Account.

* `vcenter_username` - (Required) vCenter user name for the specified SDDC. The specified user requires CloudAdmin credentials. The user does not require CloudGlobalAdmin credentials.

## Attribute Reference

* `created_at` - Date when the entity was created. The date is in ISO 6801 and UTC.

* `id` - The id of this NSX-T cloud account.

* `links` - HATEOAS of the entity.

* `org_id` - The id of the organization this entity belongs to.

* `owner` - Email of the user that owns the entity.

* `region-ids` - A set of region IDs that are enabled for this account.

* `updated_at` - Date when the entity was last updated. The date is ISO 8601 and UTC.


## Import

VMC cloud account can be imported using the id, e.g.

`$ terraform import vra_cloud_account_vmc.new_vmc 05956583-6488-4e7d-84c9-92a7b7219a15`