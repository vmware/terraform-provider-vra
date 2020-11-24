---
layout: "vra"
page_title: "VMware vRealize Automation: vra_region_enumeration_azure"
description: |-
  Provides a data lookup for region enumeration for Azure cloud account.
---

# Data Source: vra_region_enumeration_azure
## Example Usages

This is an example of how to lookup a region enumeration data source for Azure cloud account.

**Region enumeration azure data source by its id:**
```hcl
data "vra_region_enumeration_azure" "this" {
	application_id = this.applicationID
	application_key = this.applicationKey
	subscription_id = this.subscriptionID
	tenant_id = this.tenantID
 }
```

The region enumeration data source for Azure cloud account suports the following arguments:

## Argument Reference
* `application_id` - (Required) Azure Client Application ID

* `application_key` - (Required) Azure Client Application Secret Key

* `subscription_id` - (Required) Azure Subscribtion ID

* `tenant_id` - (Required) Azure Tenant ID

## Attribute Reference
* `id` - The id of the region enumeration for Azure account.

* `regions` - A set of Region names to enable provisioning on. Example: northamerica-northeast1 
