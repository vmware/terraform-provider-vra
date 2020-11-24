---
layout: "vra"
page_title: "VMware vRealize Automation: vra_region_enumeration_aws"
description: |-
  Provides a data lookup for region enumeration for AWS cloud account.
---

# Data Source: vra_region_enumeration_aws
## Example Usages

This is an example of how to lookup a region enumeration data source for AWS cloud account.

**Region enumeration AWS data source by its id:**
```hcl
data "vra_region_enumeration_aws" "this" {
  access_key = this.id
  secret_key = this.secret
}
```

The region enumeration data source for AWS cloud account suports the following arguments:

## Argument Reference
* `access_key` - (Required) Aws Access key ID. Example: ACDC55DB4MFH6ADG75KK

* `secret_key` - (Required) Aws Secret Access Key. Example: gfsScK345sGGaVdds222dasdfDDSSasdfdsa34fS

## Attribute Reference
* `id` - The id of the region enumeration for AWS account.

* `regions` - A set of Region names to enable provisioning on. Example: us-east-2, ap-northeast-1

