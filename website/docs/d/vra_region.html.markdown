---
layout: "vra"
page_title: "VMware vRealize Automation: vra_region"
description: |-
  Provides a data lookup for region data source.
---

# Data Source: vra_region
## Example Usages

This is an example of how to read a region data source.

DeprecationMessage: 'region_enumeration' is deprecated. Use 'region_enumeration_vsphere' instead.

The region data source suports the following arguments:

## Argument Reference
* `accept_self_signed_cert` - (Optional) Accept self signed certificate when connecting to vSphere. Example: false

* `dcId` - (Optional) ID of a data collector vm deployed in the on premise infrastructure. Example: d5316b00-f3b8-4895-9e9a-c4b98649c2ca

* `hostname` - (Required) Host name for the cloud account endpoint. Example: dc1-lnd.mycompany.com

* `password` - (Required) Password for the user used to authenticate with the cloud Account

* `username` - (Required) Username to authenticate with the cloud account

## Attribute Reference
* `regions` - A set of datacenter managed object reference identifiers to enable provisioning on. Example: Datacenter:datacenter-2 

