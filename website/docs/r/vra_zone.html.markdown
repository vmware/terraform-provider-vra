---
layout: "vra"
page_title: "VMware vRealize Automation: vra_zone"
sidebar_current: "docs-vra-resource-zone"
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

## Required arguments

* `name` - A human-friendly name used as an identifier in APIs that support this option.  Only one of 'filter', 'id', 'name' or 'region_id' must be specified.

* `tags` - A set of tag keys and optional values that were set on this Network Profile.

## Optional arguments

* `description` - A human-friendly description.

* `folder` - The folder relative path to the datacenter where resources are deployed to. (only applicable for vSphere cloud zones)

* `placement_policy` - The id of the region for which this zone is defined

* `region_id` - A link to the region that is associated with the storage profile.                      example:[ { "key" : "ownedBy", "value": "Rainpole" } ]
                      
* `tags_to_match` - A set of tag keys and optional values for compute resource filtering.
                   example:[ { "key" : "compliance", "value": "pci" } ]
