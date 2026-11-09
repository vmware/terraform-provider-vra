---
page_title: "VMware Aria Automation: Resource vra_deployment_resource_action"
description: A resource that can be used to run a day 2 action on a resource of a VMware Aria Automation deployment.
---

# Resource: vra_deployment_resource_action

This resource runs a day 2 action on an individual resource of a deployment, for example powering off a machine, resizing it, taking a snapshot of it, or running a custom XaaS action.

The `vra_deployment` resource runs the deployment level `Update` action when its `inputs` change. Some catalog items do not offer a deployment level `Update` action and expose day 2 actions on their individual resources instead. This resource covers that case, and it also gives access to the day 2 actions of the native infrastructure resource types.

Use the [`vra_deployment_resource_actions`](../data-sources/vra_deployment_resource_actions.md) data source to discover the actions that are available on a resource together with the inputs that each of them declares.

~> **Note:** This resource records the action that was submitted and the inputs that were used. VMware Aria Automation does not report the inputs of a completed action as the state of a resource, so this resource does not detect changes that are made outside of Terraform. The action is run again when `action_id`, `inputs`, or `triggers` change.

~> **Note:** By default, the provider waits for every action request to finish. Waiting actions submitted for the same deployment by one provider process are serialized because VMware Aria Automation accepts only one active request per deployment. With `wait_for_completion = false`, Terraform returns after submission, does not report a later terminal failure during that apply, and cannot serialize another action against the still-running request. Use that mode only when no other action will target the deployment until the request finishes. Separate concurrent Terraform runs must still be prevented by the operator or remote-state locking.

## Example Usages

This is an example of how to power off a machine of a deployment.

```hcl
resource "vra_deployment_resource_action" "power_off" {
  deployment_id = vra_deployment.this.id
  resource_name = "Cloud_vSphere_Machine_1"
  action_id     = "Cloud.vSphere.Machine.PowerOff"
}
```

This is an example of how to run an action that takes inputs. The values of the `array` and `object` inputs are supplied as JSON encoded strings, in the same way as the `inputs` of the `vra_deployment` resource.

```hcl
resource "vra_deployment_resource_action" "snapshot" {
  deployment_id = vra_deployment.this.id
  resource_id   = data.vra_deployment.this.resources[0].id
  action_id     = "Cloud.vSphere.Machine.Snapshot.Create"

  inputs = {
    name        = "before-upgrade"
    description = "Taken by Terraform"
  }

  timeouts {
    create = "60m"
    update = "60m"
  }
}
```

This is an example of a deployment whose inputs can only be changed through a day 2 action after the deployment has been created. The deployment request applies the initial values, so the action is not run when it is created. The keys that the action owns are ignored on the deployment so that the two resources do not fight over them.

```hcl
resource "vra_deployment" "this" {
  name            = var.deployment_name
  catalog_item_id = var.catalog_item_id
  project_id      = var.project_id

  inputs = {
    hostname         = var.hostname
    destination_port = var.destination_port
  }

  lifecycle {
    ignore_changes = [
      inputs["destination_port"],
    ]
  }
}

resource "vra_deployment_resource_action" "destination_port" {
  deployment_id = vra_deployment.this.id
  resource_name = "Application"
  action_id     = "Custom.UpdateApplication"

  # The deployment request has already applied this value.
  run_on_create = false

  inputs = {
    destinationPort = var.destination_port
  }
}
```

This is an example of an action that is run when the resource is destroyed.

```hcl
resource "vra_deployment_resource_action" "snapshot" {
  deployment_id = vra_deployment.this.id
  resource_name = "Cloud_vSphere_Machine_1"

  action_id = "Cloud.vSphere.Machine.Snapshot.Create"
  inputs = {
    name = "managed-by-terraform"
  }

  destroy_action_id = "Cloud.vSphere.Machine.Snapshot.Delete"
  destroy_inputs = {
    name = "managed-by-terraform"
  }
}
```

## Argument Reference

Create your resource action with the following arguments:

* `action_id` - (Required) The id of the day 2 action to run on the resource, for example `Cloud.vSphere.Machine.Snapshot.Create`. Use the `vra_deployment_resource_actions` data source to discover the available action ids.

* `deployment_id` - (Required) The id of the deployment that owns the resource. Changing this forces a new resource to be created.

* `destroy_action_id` - (Optional) The id of an action to run when this resource is destroyed. When omitted, destroying this resource only removes it from the Terraform state and does not call VMware Aria Automation.

* `destroy_inputs` - (Optional) The inputs for the action that is identified by `destroy_action_id`.

* `inputs` - (Optional) The inputs for the action. Undeclared inputs fail the apply, and values are converted using the top-level types declared by the action schema. VMware Aria Automation validates remaining schema constraints. The values of the `array` and `object` types must be supplied as JSON encoded strings.

* `resource_id` - (Optional) The id of the deployment resource to run the action on. Exactly one of `resource_id` or `resource_name` must be set. Changing this forces a new resource to be created.

* `resource_name` - (Optional) The name of the deployment resource to run the action on. The name must be unique within the deployment. Exactly one of `resource_id` or `resource_name` must be set. Changing this forces a new resource to be created.

* `run_on_create` - (Optional) Whether to run the action when this resource is created. Set this to `false` when the deployment request has already applied the desired values and the action should only run on subsequent changes. Defaults to `true`.

* `triggers` - (Optional) An arbitrary map of values that, when changed, causes the action to run again.

* `wait_for_completion` - (Optional) Whether to wait for the action request to reach a terminal state before continuing. Defaults to `true`.

## Attribute Reference

* `id` - The id of the resource action, in the format `deployment_id/resource_id/action_id`.

* `last_request_id` - The id of the most recent action request that was submitted by this resource.

* `last_request_status` - The status of the most recent action request that was submitted by this resource.

## Import

A resource action can be imported using the id, e.g.

`$ terraform import vra_deployment_resource_action.this 05956583-6488-4e7d-84c9-92a7b7219a15/2f3e4d5c-6b7a-4980-9a1b-2c3d4e5f6a7b/Cloud.vSphere.Machine.PowerOff`

~> **Note:** VMware Aria Automation does not expose the inputs used by a completed action. Import therefore cannot reconstruct `inputs` or `triggers`. Adding either to configuration after import causes Terraform to submit the action on the next apply. Review that plan before applying it.
