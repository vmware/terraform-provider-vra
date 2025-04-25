---
page_title: "VMware Aria Automation: vra_policy_iaas_resource"
description: A resource for IaaS Resource policies.
---

# Resource: vra_policy_iaas_resource

Creates an IaaS Resource policy resource to manage IaaS resources lifecycle at namespace level.

## Example Usages

The following example shows how to create an IaaS resource policy resource:

```hcl
resource "vra_policy_iaas_resource" "policy_iaas_resource" {
  name             = "terraform-iaas-resource-policy"
  description      = "IaaS Resource Policy [terraform-iaas-resource-policy] created by Terraform"
  enforcement_type = "HARD"

  failure_policy = "Fail"
  resource_rules {
    api_groups = [
      "vmoperator.vmware.com",
    ]
    api_versions = [
      "*",
    ]
    operations = [
      "CREATE",
    ]
    resources = [
      "virtualmachines",
    ]
  }
  validation_actions = [
    "Deny",
  ]
  validations {
    expression = "request.resource.resource != \"virtualmachines\""
    message    = "Virtual Machines are prohibited to be provisioned in the namespace."
  }
}
```

## Argument Reference

Create your resource with the following arguments:

* `criteria` - (Optional) The policy criteria.

* `description` - (Optional) A human-friendly description for the policy instance.

* `enforcement_type` - (Required) The type of enforcement for the policy. Supported values: `HARD`, `SOFT`.

* `exclude_resource_rules` - (Optional) Exclude Resource Rules:

  * `api_groups` - (Required) List of API groups the resources belong to.

  * `api_versions` - (Required) List of API Versions the resources belong to.

  * `operations` - (Required) List of Operations the admission hook cares about. Supported values: `CREATE`, `UPDATE`, `DELETE`.

  * `resources` - (Required) List of Resources this rule applies to.

* `failure_policy` - (Required) Failure policy to apply when the policy fails. Supported values: `Fail`, `Ignore`.

* `match_conditions` - (Optional) List of conditions that must be met for a request to be validated:

  * `expression` - (Required) Expression which will be evaluated by CEL.

  * `name` - (Required) Identifier for this match condition.

* `match_expressions` - (Optional) List of label selector requirements that must be met for an object to be validated:

  * `key` - (Required) The label key that the selector applies to.

  * `operator` - (Required) A key's relationship to a set of values.

  * `values` - (Required) An array of string values.

* `match_labels` - (Optional) Map of {key,value} pairs that must be met for an object to be validated.

* `match_policy` - (Optional) Match policy. Supported values: `Exact`, `Equivalent`.

* `name` - (Required) A human-friendly name used as an identifier for the policy instance.

* `project_criteria` - (Optional) The project based criteria. Updating this argument triggers a recreation of the resource. It cannot be specified when `project_id` is set.

* `project_id` - (Optional) The id of the project this entity belongs to. Updating this argument triggers a recreation of the resource.

* `resource_rules` - (Required) Resource Rules:

  * `api_groups` - (Required) List of API groups the resources belong to.

  * `api_versions` - (Required) List of API Versions the resources belong to.

  * `operations` - (Required) List of Operations the admission hook cares about. Supported values: `CREATE`, `UPDATE`, `DELETE`.

  * `resources` - (Required) List of Resources this rule applies to.

* `validation_actions` - (Required) List of validation actions.

* `validations` - (Required) List of CEL expressions which are used to validate admission requests:

  * `expression` - (Required) Expression which will be evaluated by CEL.

  * `message` - (Optional) Message displayed when validation fails.

  * `message_expression` - (Optional) CEL expression that evaluates to the validation failure message that is returned when this rule fails.

  * `reason` - (Optional) Machine-readable description of why this validation failed.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `created_at` - Date when the entity was created. The date is in ISO 8601 and UTC.

* `created_by` - The user the entity was created by.

* `last_updated_at` - Date when the entity was last updated. The date is ISO 8601 and UTC.

* `last_updated_by` - The user the entity was last updated by.

* `org_id` - The id of the organization this entity belongs to.

## Import

To import an existing IaaS Resource policy, use the `id` as in the following example:

`$ terraform import vra_policy_iaas_resource.policy_iaas_resource "c222fd4c-be40-43c9-a806-81ef25bdf661"`
