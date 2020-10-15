---
layout: "vra"
page_title: "VMware vRealize Automation: vra_zone"
sidebar_current: "docs-vra-datasource-vra-zone"
description: |-
  Provides a data lookup for vra_zone.
---

# Data Source: vra_zone
## Example Usages
This is an example of how to create a zone data source.

```hcl
data "vra_zone" "test-zone" {
  name = "${vra_zone.my-zone.name}"
}
```

A zone data source supports the following arguments:

## Optional arguments

* `custom_properties` - A list of key value pair of properties that will be used.

* `id` - The id of the image profile instance.  Only one of 'filter', 'id', 'name' or 'region_id' must be specified.

* `name` - A human-friendly name used as an identifier in APIs that support this option.  Only one of 'filter', 'id', 'name' or 'region_id' must be specified.

* `owner` - Email of the user that owns the entity.

## Imported attributes
* `created_at` - Date when the entity was created. The date is in ISO 8601 and UTC.

* `description` - A human-friendly description.

* `folder` - The folder relative path to the datacenter where resources are deployed to. (only applicable for vSphere cloud zones)

* `org_id` - The id of the organization this entity belongs to.

* `placement_policy` - The id of the region for which this zone is defined

* `region_id` - A link to the region that is associated with the storage profile.

* `shared_resources` - The id of the organization this entity belongs to.

* `tags` - A set of tag keys and optional values that were set on this Network Profile.
                      example:[ { "key" : "ownedBy", "value": "Rainpole" } ]
                      
* `tags_to_match` - A set of tag keys and optional values for compute resource filtering.
                   example:[ { "key" : "compliance", "value": "pci" } ]

* `updated_at` - Date when the entity was last updated. The date is ISO 8601 and UTC.
