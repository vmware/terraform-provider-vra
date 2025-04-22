---
page_title: "VMware Aria Automation: vra_policy_approval"
description: A resource for Approval policies.
---

# Resource: vra_policy_approval

Creates an Approval policy resource to request approval for catalog deployments from specified users.

## Example Usages

The following example shows how to create an approval policy resource:

```hcl
resource "vra_policy_approval" "policy_approval" {
  name             = "terraform-approval-policy"
  description      = "Approval Policy [terraform-approval-policy] created by Terraform"
  enforcement_type = "HARD"

  actions = [
    "Deployment.ChangeLease",
  ]
  approval_level = 1
  approval_mode  = "ANY_OF"
  approval_type  = "ROLE"
  approvers = [
    "ROLE:PROJECT_ADMINISTRATORS"
  ]
  auto_approval_decision = "APPROVE"
  auto_approval_expiry   = 30
}
```

## Argument Reference

Create your resource with the following arguments:

* `actions` - (Required) List of actions to trigger approval.

* `approval_level` - (Required) The level defines the order in which the policy is enforced. Level 1 approvals are applied first, followed by level 2 approvals, and so on. Valid range: `1` - `99`.

* `approval_mode` - (Required) Who must approve the request. Supported values: `ANY_OF`, `ALL_OF`.

* `approval_type` - (Required) Approval Type. Supported values: `USER`, `ROLE`.

* `approvers` - (Required) List of approvers of the policy.

* `auto_approval_decision` - (Required) Automatically approve or reject a request after the number of days specified in the Auto expiry trigger field. Supported values: `APPROVE`, `REJECT`, `NO_EXPIRY`.

* `auto_approval_expiry` - (Required) The number of days the approvers have to respond before the Auto action is triggered. Valid range: `1` - `30`.

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

To import an existing Approval policy, use the `id` as in the following example:

`$ terraform import vra_policy_approval.policy_approval "f40657aa-3089-4b80-8970-a8fa6f9f5314"`
