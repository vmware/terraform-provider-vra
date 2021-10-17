---
layout: "vra"
page_title: "VMware vRealize Automation: vra_region"
sidebar_current: "docs-vra-datasource-vra-region"
description: |-
  Provides a data lookup for vra_region.
---

# Data Source: vra\_region

This is an example of how to lookup a region data source:

**Region data source by id:**

```hcl
data "vra_region" "this" {
  id = var.vra_region_id
}
```

**Region data source by filter:**

```hcl
data "vra_region" "this" {
  filter = "name eq '${var.vra_region_name}'"
}
```

**Region data source by cloud account id and region:**

```hcl
data "vra_region" "this" {
  cloud_account_id = var.vra_cloud_account_id
  region           = var.vra_region_external_id
}
```

## Argument Reference

The following arguments are supported:

* `cloud_account_id` - (Optional) The id of the cloud account the region belongs to.

* `filter` - (Optional) Search criteria to narrow down Regions.

* `id` - (Optional) The id of the region instance.

* `region` - (Optional) The specific region associated with the cloud account. On vSphere, this is the external ID.

-> **Note:** One of `id`, `filter` or` cloud_account_id` and `region` must be specified.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `created_at` - Date when the entity was created. The date is in ISO 8601 and UTC.

* `external_region_id` - Unique identifier of region on the provider side.

* `name` - Name of region on the provider side. In vSphere, the name of the region is different from its id.

* `org_id` - The id of the organization this entity belongs to.

* `owner` - Email of the user that owns the entity.

* `updated_at` - Date when the entity was last updated. The date is ISO 8601 and UTC.
