---
layout: "vra"
page_title: "VMware vRealize Automation: vra_region_enumeration_azure"
description: |-
  Provides a data lookup for region enumeration for Azure cloud account.
---

# Data Source: vra_region_enumeration_azure
## Example Usages

This is an example of how to lookup a region enumeration data source for Azure cloud account.

**Region enumeration data source for Azure**
```hcl
data "vra_region_enumeration_azure" "this" {
  application_id  = var.application_id
  application_key = var.application_key
  subscription_id = var.subscription_id
  tenant_id       = var.tenant_id
 }
```

The region enumeration data source for Azure cloud account supports the following arguments:

## Argument Reference
* `application_id` - (Required) Azure Client Application ID

* `application_key` - (Required) Azure Client Application Secret Key

* `subscription_id` - (Required) Azure Subscription ID

* `tenant_id` - (Required) Azure Tenant ID

## Attribute Reference
* `regions` - A set of Region names to enable provisioning on. Example: `["northamerica-northeast1"]`
