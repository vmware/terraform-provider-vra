---
page_title: "VMware Aria Automation: vra_policy_day2_action"
description: A resource for Day2 Action policies.
---

# Resource: vra_policy_day2_action

Creates a Day2 Action policy resource to manage what actions are available for deployed resources.

## Example Usages

The following example shows how to create a day2 action policy resource:

```hcl
resource "vra_policy_day2_action" "policy_day2_action" {
  name             = "terraform-day2-action-policy"
  description      = "Approval Policy [terraform-day2-action-policy] created by Terraform"
  enforcement_type = "HARD"

  actions = [
    "Deployment.ChangeLease",
    "Deployment.EditDeployment"
  ]
  authorities = [
    "USER:admin",
    "GROUP:vraadamadmins@",
  ]
}
```

## Argument Reference

Create your resource with the following arguments:

* `actions` - (Optional) List of allowed actions for authority/authorities.

* `authorities` - (Required) List of authorities that will be allowed to perform certain actions.

* `criteria` - (Optional) The policy criteria.

* `description` - (Optional) A human-friendly description for the policy instance.

* `enforcement_type` - (Required) The type of enforcement for the policy. Supported values: `HARD`, `SOFT`.

* `name` - (Required) A human-friendly name used as an identifier for the policy instance.

* `project_criteria` - (Optional) The project based criteria. Updating this argument triggers a recreation of the resource. It cannot be specified when `project_id` is set.

* `project_id` - (Optional) The id of the project this entity belongs to. Updating this argument triggers a recreation of the resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `created_at` - Date when the entity was created. The date is in ISO 8601 and UTC.

* `created_by` - The user the entity was created by.

* `last_updated_at` - Date when the entity was last updated. The date is ISO 8601 and UTC.

* `last_updated_by` - The user the entity was last updated by.

* `org_id` - The id of the organization this entity belongs to.

## Import

To import an existing Day2 Action policy, use the `id` as in the following example:

`$ terraform import vra_policy_day2_action.policy_day2_action "b8c9cb7f-1faf-474d-8b9b-27ba3f1c7930"`
