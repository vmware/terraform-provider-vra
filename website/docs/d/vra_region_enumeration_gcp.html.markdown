---
layout: "vra"
page_title: "VMware vRealize Automation: vra_region_enumeration_gcp"
description: |-
  Provides a data lookup for region enumeration for GCP cloud account.
---

# Data Source: vra_region_enumeration_gcp
## Example Usages

This is an example of how to lookup a region enumeration data source for GCP cloud account.

**Region enumeration GCP data source by its id:**
```hcl
data "vra_region_enumeration_gcp" "this" {
	client_email = this.clientEmail
	project_id = this.projectID
	private_key_id = this.privateKeyID
	private_key = this.privateKey
 }
```

The region enumeration data source for GCP cloud account suports the following arguments:

## Argument Reference
* `client_email` - (Required) Client E-mail ID.

* `private_key` - (Required) GCP Private key.

* `private_key_id` - (Required) GCP Private key ID.

* `project_id` - (Required) GCP Project ID.


## Attribute Reference
* `id` - (Optional) The id of the region enumeration for GCP account.

* `regions` - A set of Region names to enable provisioning on. Example: northamerica-northeast1 

