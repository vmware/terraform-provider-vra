---
layout: "vra"
page_title: "VMware vRealize Automation: vra_region_enumeration_aws"
description: |-
  Provides a data lookup for region enumeration for AWS cloud account.
---

# Data Source: vra_region_enumeration_aws
## Example Usages

This is an example of how to lookup a region enumeration data source for AWS cloud account.

**Region enumeration data source for AWS, by the AWS account access key and secret key:**
```hcl
data "vra_region_enumeration_aws" "this" {
  access_key = var.access_key
  secret_key = var.secret_key
}
```

The region enumeration data source for AWS cloud account supports the following arguments:

## Argument Reference
* `access_key` - (Optional) Aws Access key ID.

* `secret_key` - (Required) Aws Secret Access Key.

## Attribute Reference
* `regions` - A set of Region names to enable provisioning on. Example: `["us-east-2", "ap-northeast-1"]`

