---
layout: "vra"
page_title: "VMware vRealize Automation: vra_data_collector"
description: |-
  Provides a data lookup for data collector data source.
---

# Data Source: vra_data_collector
## Example Usages

This is an example of how to lookup a data collector.

**Data collector data source by its name:**
```hcl
data "vra_data_collector" "this" {
  name = var.data_collector_name
}
```
The data collector data source supports the following arguments:

## Argument Reference
* `name` - (Required) Data collector name. Example: `Datacollector1`

## Attribute Reference
* `hostname` - Data collector host name. Example: `dc1-lnd.mycompany.com`

* `ip_address` - IPv4 Address of the data collector VM. Example: `10.0.0.1`

* `status` - Current status of the data collector. Example: `ACTIVE`, `INACTIVE`

