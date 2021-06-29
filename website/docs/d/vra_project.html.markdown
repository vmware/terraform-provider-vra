---
layout: "vra"
page_title: "VMware vRealize Automation: vra_project"
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

## Argument Reference

* `administrators` - (Optional) List of administrator users associated with the project. Only administrators can manage project's configuration. 
Deprecated, to specify the type of principal, please refer `administrator_roles`.

* `administrator_roles` - (Optional) Administrator users or groups associated with the project. Only administrators can manage project's configuration. 

* `constraints` - (Optional) List of storage, network and extensibility constraints to be applied when provisioning through this project.

* `custom_properties` - (Optional) The project custom properties which are added to all requests in this project.

* `description` - (Optional) A human-friendly description.

* `id` - (Optional) The id of the image profile instance.

* `machine_naming_template` - (Optional) The naming template to be used for resources provisioned in this project.

* `members` - (Optional) List of member users associated with the project. Deprecated, to specify the type of principal, please refer `member_roles`.

* `member_roles` - (Optional) Member users or groups associated with the project. 

* `name` - (Optional) A human-friendly name used as an identifier in APIs that support this option.

* `operation_timeout` - (Optional) The timeout that should be used for Blueprint operations and Provisioning tasks. The timeout is in seconds.

* `shared_resources` - (Optional) Specifies whether the resources in this projects are shared or not. If not set default will be used.

* `zone_assignments` - (Optional) List of configurations for zone assignment to a project.

* `shared_resources` - (Optional) The id of the organization this entity belongs to.

* `viewers` - (Optional) List of viewer users associated with the project. Deprecated, to specify the type of principal, please refer `viewer_roles`.

* `viewer_roles` - (Optional) Viewer users or groups associated with the project. 
