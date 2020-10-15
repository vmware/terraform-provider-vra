---
layout: "vra"
page_title: "VMware vRealize Automation: vra_image_profile"
sidebar_current: "docs-vra-datasource-vra-image_profile"
description: |-
  Provides a data lookup for vra_image_profile.
---

# Data Source: vra_image_profile
## Example Usages
This is an example of how to read an image profile as data source.

**Image profile data source by its id:**

```hcl
data "vra_image_profile" "this" {
  filter = "name eq 'foobar'"
}
```

**Vra image profile data source filter by region id:**

```hcl
data "vra_image_profile" "this" {
  region_id = vra_image_profile.this.region_id
}
```

An image profile data source supports the following arguments:

## Optional arguments

* `description` - A human-friendly description.

* `filter` - Filter query string that is supported by vRA multi-cloud IaaS API. Example: regionId eq '<regionId>' and cloudAccountId eq '<cloudAccountId>'. Only one of 'filter', 'id', 'name' or 'region_id' must be specified.

* `id` - The id of the image profile instance.  Only one of 'filter', 'id', 'name' or 'region_id' must be specified.

* `image_mapping` - Image mapping defined for the corresponding region.

* `name` - A human-friendly name used as an identifier in APIs that support this option.  Only one of 'filter', 'id', 'name' or 'region_id' must be specified.

* `region_id` - The id of the region for which this profile is defined as in vRealize Automation(vRA).  Only one of 'filter', 'id', 'name' or 'region_id' must be specified.

## Imported attributes

* `created_at` - Date when the entity was created. The date is in ISO 6801 and UTC.

* `external_region_id` - The external regionId of the resource. 

* `owner` - Email of the user that owns the entity.

* `updated_at` - Date when the entity was last updated. The date is ISO 8601 and UTC.
