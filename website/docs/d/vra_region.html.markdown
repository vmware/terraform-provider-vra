---
layout: "vra"
page_title: "VMware vRealize Automation: vra_region"
description: |-
  Provides a data lookup for region data source.
---

# Data Source: vra_region
## Example Usages

This is an example of how to read a region data source.

**Region data source by cloud account Id and region:**
```hcl
data "vra_region" "this" {
  cloud_account_id = var.vra_cloud_account_id
  region           = "us-east-1"
}
```

The region data source supports the following arguments:

## Argument Reference
* `cloud_account_id` - (Required) The Cloud Account ID.

* `region` - (Required) The specific region associated with the cloud account.

## Attribute Reference
* `id` - The ID of the given region within the cloud account.
