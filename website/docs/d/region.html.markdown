---
layout: "vra"
page_title: "VMware vRealize Automation: vra_region"
sidebar_current: "docs-vra-datasource-vra-region"
description: |-
  Provides a data lookup for vra_region.
---

# Data Source: vra\_region

This is an example of how to lookup a region data source.

**Region data source**
```hcl
data "vra_region" "region_1" {
  cloud_account_id = "cloud_account_id"
  region           = "us-east-1"
}
```

```hcl
data "vra_region" "eastus" {
  cloud_account_id = "cloud_account_id"
  region           = "eastus"
}
```

The region data source supports the following arguments:

## Argument Reference
* `cloud_account_id` - (Optional) The id of the cloud account this region belongs to. Example: 9e49

* `filter` - (Optional) The id of this resource instance.

* `id` - (Optional) The id of this resource instance. Example: 9e49

* `region` - (Optional) Name of region on the provider side. In vSphere, the name of the region is different from its id. Example: us-west


## Attribute Reference
* `created_at` - Date when the entity was created. The date is in ISO 8601 and UTC. Example: 2012-09-27

* `external_region_id` - Unique identifier of region on the provider side. Example: us-west

* `name` - Name of region on the provider side. In vSphere, the name of the region is different from its id. Example: us-west

* `org_id` - The id of the organization this entity belongs to. Example: 9e49

* `owner` - Email of the user that owns the entity. Example: csp@vmware.com

* `updated_at` - Date when the entity was last updated. The date is ISO 8601 and UTC. Example: 2012-09-27

