---
layout: "vra"
page_title: "VMware vRealize Automation: vra_project"
description: |-
  Provides a VMware vRA vra_project resource.
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

  # Deprecated, please use administrator_roles instead.
  administrators = ["jason@vra.local"]
  
  administrator_roles {
    email = "jason@vra.local"
    type = "user"
  }

  administrator_roles {
    email = "jason-group@vra.local"
    type = "group"
  }

  # Deprecated, please use member_roles instead.
  members = ["tony@vra.local"]

  member_roles {
    email = "tony@vra.local"
    type = "user"
  }

  member_roles {
    email = "tony-group@vra.local"
    type = "group"
  }

  # Deprecated, please use viewer_roles instead
  viewers = ["shauna@vra.local"]

  viewer_roles {
    email = "shauna@vra.local"
    type = "user"
  }

  viewer_roles {
    email = "shauna-group@vra.local"
    type = "group"
  }

  operation_timeout = 6000

  machine_naming_template = "$${resource.name}-$${####}"

  constraints {
    extensibility {
      expression = "foo:bar"
      mandatory  = false
    }
    extensibility {
      expression = "environment:Test"
      mandatory  = true
    }

    network {
      expression = "foo:bar"
      mandatory  = false
    }
    network {
      expression = "environment:Test"
      mandatory  = true
    }

    storage {
      expression = "foo:bar"
      mandatory  = false
    }
    storage {
      expression = "environment:Test"
      mandatory  = true
    }
  }
}
```

A project resource supports the following arguments:

## Argument Reference

* `administrators` - (Optional) List of administrator users associated with the project. Only administrators can manage project's configuration. 
Deprecated, specify the type of principal, please refer `administrator_roles`.

* `administrator_roles` - (Optional) Administrator users or groups associated with the project. Only administrators can manage project's configuration. 

* `constraints` - (Optional) List of storage, network and extensibility constraints to be applied when provisioning through this project.

* `custom_properties` - (Optional) The project custom properties which are added to all requests in this project.

* `description` - (Optional) A human-friendly description.

* `machine_naming_template` - (Optional) The naming template to be used for resources provisioned in this project.

* `members` - (Optional) List of member users associated with the project. Deprecated, specify the type of principal, please refer `member_roles`.

* `member_roles` - (Optional) Member users or groups associated with the project. 

* `name` - (Required) A human-friendly name used as an identifier in APIs that support this option.

* `operation_timeout` - (Optional) The timeout that should be used for Blueprint operations and Provisioning tasks. The timeout is in seconds.

* `shared_resources` - (Optional) Specifies whether the resources in this projects are shared or not. If not set default will be used.

* `viewers` - (Optional) List of viewer users associated with the project. Deprecated, specify the type of principal, please refer `viewer_roles`.

* `viewer_roles` - (Optional) Viewer users or groups associated with the project. 

* `zone_assignments` - (Optional) List of configurations for zone assignment to a project.