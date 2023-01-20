---
layout: "vra"
page_title: "VMware vRealize Automation: vra_region_enumeration_gcp"
description: |-
  Provides a data lookup for region enumeration for GCP cloud account.
---

# Data Source: vra_region_enumeration_gcp
## Example Usages

This is an example of how to lookup a region enumeration data source for GCP cloud account.

**Region enumeration data source for GCP**
```hcl
data "vra_region_enumeration_gcp" "this" {
  client_email   = var.client_email
  private_key_id = var.private_key_id
  private_key    = var.private_key
  project_id     = var.project_id
 }
```

The region enumeration data source for GCP cloud account supports the following arguments:

## Argument Reference
* `client_email` - (Required) Client E-mail ID.

* `private_key` - (Required) GCP Private key.

* `private_key_id` - (Required) GCP Private key ID.

* `project_id` - (Required) GCP Project ID.

## Attribute Reference
* `regions` - A set of Region names to enable provisioning on. Example: `["northamerica-northeast1"]`

