---
layout: "vra"
page_title: "VMware vRealize Automation: vra_flavor_profile"
description: |-
  Provides a data lookup for vra_flavor_profile.
---

# Resource: vra_flavor_profile
## Example Usages
This is an example of how to create a flavor profile resource.

**Flavor profile:**

```hcl
resource "vra_flavor_profile" "my-flavor-profile" {
	name = "AWS"
	description = "my flavor"
	region_id = "${data.vra_region.us-east-1-region.id}"
	flavor_mapping {
		name = "small"
		instance_type = "t2.small"
	}
	flavor_mapping {
		name = "medium"
		instance_type = "t2.medium"
	}
}
```

An flavor profile resource supports the following arguments:

## Argument Reference

* `description` - (Optional) A human-friendly description.

* `name` - (Required) A human-friendly name used as an identifier in APIs that support this option.

* `region_id` - (Required) The id of the region for which this profile is defined as in vRealize Automation(vRA).

* `flavor_mapping` - (Optional) Map between global fabric flavor keys and fabric flavor descriptions.

## Attribute Reference

* * `cloud_account_id` - The ID of the cloud account this flavor profile belongs to.

* `created_at` - Date when  entity was created. Date and time format is ISO 8601 and UTC.

* `external_region_id` - The ID of the external region that is associated with the flavor profile.

* `links` - HATEOAS of entity.

* `org_id` - ID of organization that entity belongs to.

* `owner` - Email of entity owner.

* `updated_at` - Date when the entity was last updated. Date and time format is ISO 8601 and UTC.