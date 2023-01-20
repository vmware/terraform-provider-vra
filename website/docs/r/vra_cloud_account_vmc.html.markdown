---
layout: "vra"
page_title: "VMware vRealize Automation: vra_cloud_account_vmc"
description: |-
    Creates a vra_cloud_account_vmc resource.
---

# Resource: vra\_cloud\_account\_vmc

Creates a VMware vRealize Automation VMC cloud account resource.

## Example Usages

**Create VMC cloud account:**

The following example shows how to create a VMC cloud account resource.

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

Create your VMC cloud account resource with the following arguments:

* `accept_self_signed_cert` - (Optional) Accept self-signed certificate when connecting to the cloud account.

* `api_token` - (Required) VMC API access key.

* `dc_id` - (Optional) Identifier of a data collector VM deployed in the on premise infrastructure. Refer to the data-collector API to create or list data collector.

* `description` - (Optional) Human-friendly description.

* `name` - (Optional) Human-friendly name used as an identifier in APIs that support this option.

* `nsx_hostname` - (Required) IP address of the NSX Manager server in the specified SDDC / FQDN.

* `regions` - (Optional) Set of region names enabled for the cloud account.

* `sddc_name` - (Required) Identifier of the on-premise SDDC to be used by the cloud account. Note that NSX-V SDDCs are not supported.

* `tags` - (Optional) Set of tag keys and values to apply to the cloud account.
Example: `[ { "key" : "vmware", "value": "provider" } ]`

* `vcenter_hostname` - (Required) IP address or FQDN of the vCenter Server in the specified SDDC. The cloud proxy belongs on this vCenter.

* `vcenter_password` - (Required) Password used to authenticate to the cloud Account.

* `vcenter_username` - (Required) vCenter user name for the specified SDDC. The user requires CloudAdmin credentials. The user does not require CloudGlobalAdmin credentials.

## Attribute Reference

* `created_at` - Date when entity was created. Date and time format is ISO 8601 and UTC.

* `id` - ID of the VMC cloud account.

* `links` - HATEOAS of entity.

* `org_id` - ID of organization that entity belongs to.

* `owner` - Email of entity owner.

* `updated_at` - Date when the entity was last updated. Date and time format is ISO 8601 and UTC.


## Import

To import the VMC cloud account, use the ID as in the following example:

`$ terraform import vra_cloud_account_vmc.new_vmc 05956583-6488-4e7d-84c9-92a7b7219a15`
