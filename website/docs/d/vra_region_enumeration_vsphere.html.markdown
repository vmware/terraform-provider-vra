---
layout: "vra"
page_title: "VMware vRealize Automation: vra_region_enumeration_vsphere"
description: |-
  Provides a data lookup for region enumeration for vSphere cloud account.
---

# Data Source: vra_region_enumeration_vsphere
## Example Usages

This is an example of how to lookup a region enumeration data source for vSphere cloud account.

**Region enumeration vSphere data source**
```hcl
data "vra_data_collector" "dc" {
		name = dc.dcname
}

data "vra_region_enumeration_vsphere" "this" {
	  username    = this.username
	  password    = this.password
	  hostname    = this.hostname
	  dcid        = data.vra_data_collector.dc.id
}
```

The region enumeration data source for vSphere cloud account suports the following arguments:

## Argument Reference
* `accept_self_signed_cert` - (Optional) Accept self signed certificate when connecting to vSphere. Example: false

* `dcId` - (Required) ID of a data collector vm deployed in the on premise infrastructure. Example: d5316b00-f3b8-4895-9e9a-c4b98649c2ca

* `hostname` - (Required) Host name for the cloud account endpoint. Example: dc1-lnd.mycompany.com

* `password` - (Required) Password for the user used to authenticate with the cloud Account

* `username` - (Required) Username to authenticate with the cloud account

## Attribute Reference
* `id` - The ID of the region enumeration for vSphere account.

* `regions` - A set of datacenter managed object reference identifiers to enable provisioning on. Example: Datacenter:datacenter-2 

