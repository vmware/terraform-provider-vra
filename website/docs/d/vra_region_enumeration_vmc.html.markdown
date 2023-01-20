---
layout: "vra"
page_title: "VMware vRealize Automation: vra_region_enumeration_vmc"
description: |-
  Provides a data lookup for region enumeration for VMC cloud account.
---

# Data Source: vra_region_enumeration_vmc
## Example Usages

This is an example of how to lookup a region enumeration data source for VMC cloud account.

**Region enumeration data source for VMC**
```hcl
data "vra_region_enumeration_vmc" "this" {
  accept_self_signed_cert = true

  dc_id = var.vra_data_collector_id

  api_token = var.api_token
  sddc_name = var.sddc_name
  nsx_hostname = var.nsx_hostname

  vcenter_hostname = var.vcenter_hostname
  vcenter_password = var.vcenter_password
  vcenter_username = var.vcenter_username
}
```

The region enumeration data source for VMC cloud account supports the following arguments:

## Argument Reference
* `accept_self_signed_cert` - (Optional) Accept self signed certificate when connecting to vSphere. Example: `false`

* `api_token` - (Required) API Token for the cloud account endpoint.

* `dc_id` - (Optional) ID of a data collector vm deployed in the on premise infrastructure. Refer to the data-collector API to create or list data collectors.

* `nsx_hostname` - (Required) The IP address of the NSX Manager server in the specified SDDC / FQDN.

* `sddc_name` - (Required) Identifier of the on-premise SDDC to be used by this cloud account.

* `vcenter_hostname` - (Required) The IP address or FQDN of the vCenter Server in the specified SDDC. The cloud proxy belongs on this vCenter.

* `vcenter_password` - (Required) Password for the user used to authenticate with the cloud Account

* `vcenter_username` - (Required) vCenter user name for the specified SDDC.The specified user requires CloudAdmin credentials. The user does not require CloudGlobalAdmin credentials.

## Attribute Reference
* `regions` - A set of Region names to enable provisioning on. Example: `["northamerica-northeast1"]`
