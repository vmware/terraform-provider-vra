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
