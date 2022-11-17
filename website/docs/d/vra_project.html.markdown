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

* `administrators` - (Optional) A list of administrator users associated with the project. Only administrators can manage project's configuration. 

  > **Note**:  Deprecated - please use `administrator_roles` instead.

* `administrator_roles` - (Optional) Administrator users or groups associated with the project. Only administrators can manage project's configuration. 

* `constraints` - (Optional) A list of storage, network and extensibility constraints to be applied when provisioning through this project.

* `custom_properties` - (Optional) The project custom properties which are added to all requests in this project.

* `description` - (Optional) A human-friendly description.

* `id` - (Optional) The id of the image profile instance.

* `machine_naming_template` - (Optional) The naming template to be used for resources provisioned in this project.

* `members` - (Optional) A list of member users associated with the project. 
  
  > **Note**:  Deprecated - please use `member_roles` instead.

* `member_roles` - (Optional) Member users or groups associated with the project. 

* `name` - (Optional) A human-friendly name used as an identifier in APIs that support this option.

* `operation_timeout` - (Optional) The timeout that should be used for Blueprint operations and Provisioning tasks. The timeout is in seconds.

* `placement_policy` - (Optional) The placement policy that will be applied when selecting a cloud zone for provisioning. Must be one of `DEFAULT` or `SPREAD`.

* `shared_resources` - (Optional) Specifies whether the resources in this projects are shared or not. If not set default will be used.

* `zone_assignments` - (Optional) A list of configurations for zone assignment to a project.

* `shared_resources` - (Optional) The id of the organization this entity belongs to.

* `supervisor_roles` - (Optional) Supervisor users or groups associated with the project.d

* `viewers` - (Optional) A list of viewer users associated with the project. 

  > **Note**:  Deprecated - please use `viewer_roles` instead.

* `viewer_roles` - (Optional) Viewer users or groups associated with the project. 
