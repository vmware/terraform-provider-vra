---
layout: "vra"
page_title: "VMware vRealize Automation: vra_region_enumeration_vmc"
description: |-
  Provides a data lookup for region enumeration for VMC cloud account.
---

# Data Source: vra_region_enumeration_vmc
## Example Usages

This is an example of how to lookup a region enumeration data source for VMC cloud account.

**Region enumeration VMC data source**
```hcl
data "vra_data_collector" "dc" {
	name = dc.dcname
}

data "vra_region_enumeration_vmc" "this" {
	api_token = this.apiToken
	sddc_name = this.sddcName
	vcenter_username  = this.vCenterUserName
	vcenter_password  = this.vCenterPassword
	vcenter_hostname  = this.vCenterHostName
	dc_id = data.vra_data_collector.dc.id
	accept_self_signed_cert = true
}
```

The region enumeration data source for VMC cloud account suports the following arguments:

## Argument Reference
* `accept_self_signed_cert` - (Optional) Accept self signed certificate when connecting to vSphere. Example: false

* `api_token` - (Required) Host name for the cloud account endpoint. Example: dc1-lnd.mycompany.com

* `dc_id` - (Required) ID of a data collector vm deployed in the on premise infrastructure. Refer to the data-collector API to create or list data collectors.

* `sddc_name` - (Required) Identifier of the on-premise SDDC to be used by this cloud account.

* `vcenter_hostname` - (Required) The IP address or FQDN of the vCenter Server in the specified SDDC. The cloud proxy belongs on this vCenter.

* `vcenter_password` - (Required) Password for the user used to authenticate with the cloud Account

* `vcenter_username` - (Required) vCenter user name for the specified SDDC.The specified user requires CloudAdmin credentials. The user does not require CloudGlobalAdmin credentials.

## Attribute Reference
* `id` - The id of the region enumeration for VMC account.

* `regions` - A set of Region names to enable provisioning on. Example: northamerica-northeast1 