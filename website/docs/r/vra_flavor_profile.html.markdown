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
resource "vra_flavor_profile" "my-aws-flavor-profile" {
	name = "AWS"
	description = "My AWS flavor"
	region_id = "${data.vra_region.aws.id}"
	flavor_mapping {
		name = "small"
		instance_type = "t2.small"
	}
	flavor_mapping {
		name = "medium"
		instance_type = "t2.medium"
	}
}

resource "vra_flavor_profile" "my-vsphere-flavor-profile" {
	name = "vSphere"
	description = "My vSphere flavor"
	region_id = "${data.vra_region.vsphere.id}"
	flavor_mapping {
		name = "small"
		cpu_count = 2
		memory = 2048
	}
	flavor_mapping {
		name = "medium"
		cpu_count = 4
		memory = 4096
	}
}
```

An flavor profile resource supports the following arguments:

## Argument Reference

* `description` - (Optional) A human-friendly description.

* `name` - (Required) A human-friendly name used as an identifier in APIs that support this option.

* `region_id` - (Required) The id of the region for which this profile is defined.

* `flavor_mapping` - (Optional) A list of the flavor mappings defined for the corresponding cloud end-point region.
	* `name` - (Required) The name of the flavor mapping.
	* `instance_type` - (Optional) The value of the instance type in the corresponding cloud. Mandatory for public clouds. Only `instance_type` or `cpu_count`/`memory` must be specified.
	* `cpu_count` - (Optional) Number of CPU cores. Mandatory for private clouds such as vSphere. Only `instance_type` or `cpu_count`/`memory` must be specified.
	* `memory` - (Optional) Total amount of memory (in megabytes). Mandatory for private clouds such as vSphere. Only `instance_type` or `cpu_count`/`memory` must be specified.

## Attribute Reference

* `cloud_account_id` - Id of the cloud account this flavor profile belongs to.

* `created_at` - Date when the entity was created. The date is in ISO 8601 and UTC.

* `external_region_id` - The id of the region for which this profile is defined.

* `links` - HATEOAS of entity.

* `org_id` - The id of the organization this entity belongs to.

* `owner` - Email of the user that owns the entity.

* `updated_at` - Date when the entity was last updated. The date is ISO 8601 and UTC.