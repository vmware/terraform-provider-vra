---
layout: "vra"
page_title: "VMware vRealize Automation: vra_project"
description: |-
  Provides a VMware vRealize Automation vra_project resource.
---
# Resource: vra\_project

## Example Usages

This is an example of how to create a project resource.

```hcl
resource "vra_project" "this" {
  name        = var.project_name
  description = "terraform test project"

  zone_assignments {
    zone_id          = data.vra_zone.this.id
    priority         = 1
    max_instances    = 2
    cpu_limit        = 1024
    memory_limit_mb  = 8192
    storage_limit_gb = 65536
  }

  custom_properties = {
    "foo": "bar",
    "foo2": "bar2"
  }

  shared_resources = false

  administrator_roles {
    email = "jason@vra.local"
    type  = "user"
  }

  administrator_roles {
    email = "jason-group@vra.local"
    type  = "group"
  }

  member_roles {
    email = "tony@vra.local"
    type  = "user"
  }

  member_roles {
    email = "tony-group@vra.local"
    type  = "group"
  }

  supervisor_roles {
    email = "ethan@vra.local"
    type  = "user"
  }

  supervisor_roles {
    email = "ethan-group@vra.local"
    type  = "group"
  }

  viewer_roles {
    email = "shauna@vra.local"
    type  = "user"
  }

  viewer_roles {
    email = "shauna-group@vra.local"
    type  = "group"
  }

  operation_timeout = 6000

  machine_naming_template = "$${resource.name}-$${####}"

  placement_policy = "SPREAD"

  constraints {
    extensibility {
      expression = "foo:bar"
      mandatory  = false
    }
    extensibility {
      expression = "environment:test"
      mandatory  = true
    }

    network {
      expression = "foo:bar"
      mandatory  = false
    }
    network {
      expression = "environment:test"
      mandatory  = true
    }

    storage {
      expression = "foo:bar"
      mandatory  = false
    }
    storage {
      expression = "environment:test"
      mandatory  = true
    }
  }
}
```

A project resource supports the following arguments:

## Argument Reference

* `administrators` - (Optional) A list of administrator users associated with the project. Only administrators can manage project's configuration.

  > **Note**:  Deprecated - please use `administrator_roles` instead.

* `administrator_roles` - (Optional) Administrator users or groups associated with the project. Only administrators can manage project's configuration.

* `constraints` - (Optional) A list of storage, network, and extensibility constraints to be applied when provisioning through this project.

* `custom_properties` - (Optional) The project custom properties which are added to all requests in this project.

* `description` - (Optional) A human-friendly description.

* `machine_naming_template` - (Optional) The naming template to be used for resources provisioned in this project.

* `members` - (Optional) A list of member users associated with the project.

  > **Note**:  Deprecated - please use `member_roles` instead.

* `member_roles` - (Optional) Member users or groups associated with the project.

* `name` - (Required) A human-friendly name used as an identifier in APIs that support this option.

* `operation_timeout` - (Optional) The timeout that should be used for cloud template operations and provisioning tasks. The timeout is measured in seconds.

* `placement_policy` - (Optional) The placement policy that will be applied when selecting a cloud zone for provisioning. Must be one of `DEFAULT` or `SPREAD`.

* `shared_resources` - (Optional) Specifies whether the resources in this projects are shared or not. If not set default will be used.

* `supervisor_roles` - (Optional) Supervisor users or groups associated with the project.

* `viewers` - (Optional) A list of viewer users associated with the project.

  > **Note**:  Deprecated - please use `viewer_roles` instead.

* `viewer_roles` - (Optional) Viewer users or groups associated with the project.

* `zone_assignments` - (Optional) A list of configurations for zone assignment to a project.

**Due to the design of the vRealize Automation IaaS API to update a project, it's not able to add and remove user or group at the same time. Please execute `terraform apply` twice.**

Initially, we have `jason` and `tony` configured as administrator. The initial the configuration:

```hcl
  administrator_roles {
    email = "jason@vra.local"
    type = "user"
  }

  administrator_roles {
    email = "tony@vra.local"
    type = "user"
  }
```

Next, we want to add `bob` as a new administrator and remove `jason`. The modified configuration:

```hcl
  administrator_roles {
    email = "bob@vra.local"
    type = "user"
  }

  administrator_roles {
    email = "tony@vra.local"
    type = "user"
  }
```

To complete the whole operation, it requires running `terraform apply` twice.

