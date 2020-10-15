---
layout: "vra"
page_title: "VMware vRealize Automation: vra_project"
sidebar_current: "docs-vra-datasource-vra-project"
description: |-
  Provides a data lookup for vra_project.
---

# Data Source: vra_project
## Example Usages
This is an example of how to create a project data source.

**Project data source by its id:**

```hcl
data "vra_project" "this" {
  id = "${vra_project.my-project.id}"
}
```

**Project data source filter by name:**

```hcl
data "vra_project" "test-project" {
  name = "${vra_project.my-project.name}"
}
```

A project data source supports the following arguments:

## Optional arguments

* `administrators` - List of administrator users associated with the project. Only administrators can manage project's configuration.

* `constraints` - List of storage, network and extensibility constraints to be applied when provisioning through this project.

* `description` - A human-friendly description.

* `id` - The id of the image profile instance.  Only one of 'filter', 'id', 'name' or 'region_id' must be specified.

* `machine_naming_template` - The naming template to be used for resources provisioned in this project.

* `members` - List of member users associated with the project.

* `name` - A human-friendly name used as an identifier in APIs that support this option.  Only one of 'filter', 'id', 'name' or 'region_id' must be specified.

* `operation_timeout` - The timeout that should be used for Blueprint operations and Provisioning tasks. The timeout is in seconds.

* `shared_resources` - Specifies whether the resources in this projects are shared or not. If not set default will be used.

* `zone_assignments` - List of configurations for zone assignment to a project.

* `shared_resources` - The id of the organization this entity belongs to.

* `viewer` - List of viewer users associated with the project.
