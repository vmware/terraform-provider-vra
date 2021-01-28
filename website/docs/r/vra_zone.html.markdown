---
layout: "vra"
page_title: "VMware vRealize Automation: vra_zone"
description: |-
  Provides a VMware vRA vra_zone resource.
---

# vra_zone
# Resource: vra_zone
## Example Usages
This is an example of how to create a zone resource.

```hcl
resource "vra_zone" "this" {
  name        = "tf-vra-zone"
  description = "my terraform test cloud zone"
  region_id   = data.vra_region.this.id

  tags {
    key   = "my-tf-key"
    value = "my-tf-value"
  }

  tags {
    key   = "tf-foo"
    value = "tf-bar"
  }
}
```

A zone profile resource supports the following arguments:

## Argument Reference

* `description` - (Optional) A human-friendly description.

* `folder` - (Optional) The folder relative path to the datacenter where resources are deployed to. (only applicable for vSphere cloud zones)

* `name` - (Required) A human-friendly name used as an identifier in APIs that support this option.

* `placement_policy` - (Optional) The id of the region for which this zone is defined. Valid values are: `DEFAULT`, `SPREAD`, `BINPACK`. Default is `DEFAULT`.

* `region_id` - (Required) A link to the region that is associated with the zone.                      example:[ { "key" : "ownedBy", "value": "Rainpole" } ]

* `tags` - (Optional) A set of tag keys and optional values that were set on this zone.

* `tags_to_match` - (Optional) A set of tag keys and optional values for compute resource filtering.
                   example:[ { "key" : "compliance", "value": "pci" } ]

## Attribute Reference

* `cloud_account_id` - The ID of the cloud account this zone belongs to.

* `created_at` - Date when  entity was created. Date and time format is ISO 8601 and UTC.

* `custom_properties` - A list of key value pair of properties related to the zone.

* `external_region_id` - The ID of the external region that is associated with the zone.

* `links` - HATEOAS of entity.

* `org_id` - ID of organization that entity belongs to.

* `owner` - Email of entity owner.

* `updated_at` - Date when the entity was last updated. Date and time format is ISO 8601 and UTC.