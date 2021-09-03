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

**Region data source by filter:**
```hcl
data "vra_region" "this" {
  filter = "externalRegionId eq '${var.external_region_id}' and cloudAccountId eq '${vra_cloud_account.this.id}'"
}
```

**Region data source by filter and name:**
For setting up vra_zone, vra_image_profile, vra_network_profile, etc for a VMC cloud account, use the following example.
```hcl
data "vra_region" "this" {
  filter = "externalRegionId eq '${var.external_region_id}' and cloudAccountId ne '${vra_cloud_account_vmc.this.id}'"
  name   = var.region_name
}
```

The region data source supports the following arguments:

## Argument Reference
* `cloud_account_id` - (Optional) The id of the cloud account the region belongs to.

* `filter` - (Optional) Search criteria to narrow down Regions.

* `id` - (Optional) The id of the region instance,

* `region` - (Optional) The specific region associated with the cloud account.

* `name` - (Optional) Name of region on the provider side. In vSphere, the name of the region is different from its id.

## Attribute Reference
* `created_at` - Date when the entity was created. The date is in ISO 8601 and UTC.

* `org_id` - The id of the organization this entity belongs to.

* `owner` - Email of the user that owns the entity.

* `external_region_id` - Unique identifier of region on the provider side.

* `id` - The id of this resource instance.

* `updated_at` - Date when the entity was last updated. The date is ISO 8601 and UTC.
