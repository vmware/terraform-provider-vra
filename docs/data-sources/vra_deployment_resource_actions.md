---
page_title: "VMware Aria Automation: Data Source vra_deployment_resource_actions"
description: A data source that lists the day 2 actions that are available on a resource of a VMware Aria Automation deployment.
---

# Data Source: vra_deployment_resource_actions

This data source lists the day 2 actions that are available on an individual resource of a deployment, together with the schema of the inputs that each of them declares.

Use it to discover the values for the `action_id` and `inputs` arguments of the [`vra_deployment_resource_action`](../resources/vra_deployment_resource_action.md) resource.

## Example Usages

This is an example of how to list the actions that are available on a resource of a deployment.

```hcl
data "vra_deployment_resource_actions" "this" {
  deployment_id = vra_deployment.this.id
  resource_name = "Cloud_vSphere_Machine_1"
}

output "action_ids" {
  value = [for action in data.vra_deployment_resource_actions.this.actions : action.id]
}
```

This is an example of how to list only the actions that can be run against the resource in its current state, and inspect the inputs that one of them declares.

```hcl
data "vra_deployment_resource_actions" "this" {
  deployment_id = vra_deployment.this.id
  resource_name = "Cloud_vSphere_Machine_1"
  valid_only    = true
}

output "snapshot_action_schema" {
  value = one([
    for action in data.vra_deployment_resource_actions.this.actions :
    jsondecode(action.schema_json) if action.id == "Cloud.vSphere.Machine.Snapshot.Create"
  ])
}
```

## Argument Reference

* `deployment_id` - (Required) The id of the deployment that owns the resource.

* `resource_id` - (Optional) The id of the deployment resource. Exactly one of `resource_id` or `resource_name` must be set.

* `resource_name` - (Optional) The name of the deployment resource. The name must be unique within the deployment. Exactly one of `resource_id` or `resource_name` must be set.

* `valid_only` - (Optional) Whether to return only the actions that can be run against the resource in its current state. Defaults to `false`.

## Attribute Reference

* `actions` - The day 2 actions that are available on the deployment resource.

  * `description` - The description of the action.

  * `display_name` - The display name of the action.

  * `id` - The id of the action, to be used as the `action_id` of a `vra_deployment_resource_action` resource.

  * `name` - The name of the action.

  * `schema_json` - The schema of the inputs of the action, in the encoded JSON string format.

  * `valid` - Whether the action can be run against the resource in its current state.

* `resource_id` - The id of the deployment resource.
