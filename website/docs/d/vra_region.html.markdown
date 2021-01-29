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
* `cloud_account_id` - (Optional) The Cloud Account ID.

* `filter` - (Optional) Search criteria to narrow down Images.

* `id` - (Optional) The ID of the region to find.

* `region` - (Optional) The specific region associated with the cloud account.

* `name` - (Optional) Name of the region from the associated vCenter.

## Attribute Reference
* `created_at` - Date when entity was created. Date and time format is ISO 8601 and UTC.

* `org_id` - ID of organization that entity belongs to.

* `owner` - Email of entity owner.

* `external_region_id` - External regionId of region.

* `id` - The ID of the given region within the cloud account.

* `updated_at` - Date when entity was updated. Date and time format is ISO 8601 and UTC.
