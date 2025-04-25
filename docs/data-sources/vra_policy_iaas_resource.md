---
page_title: "VMware Aria Automation: vra_policy_iaas_resource"
description: A data source for IaaS Resource policies.
---

# Data Source: vra_policy_iaas_resource

The following examples shows how to lookup for an IaaS resource policy:

**IaaS Resource policy data source by its id:**

```hcl
data "vra_policy_iaas_resource" "this" {
  id = var.vra_iaas_resource_policy_id
}
```

**IaaS Resource policy data source by name search:**

```hcl
data "vra_policy_iaas_resource" "this" {
  search = var.vra_iaas_resource_policy_search_name
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Optional) The id of the policy instance.

* `search` - (Optional) Search criteria to narrow down the policy instance.

-> **Note:** One of `id` or `search` must be specified.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `created_at` - Date when the entity was created. The date is in ISO 8601 and UTC.

* `created_by` - The user the entity was created by.

* `criteria` - The policy criteria.

* `description` - A human-friendly description for the policy instance.

* `enforcement_type` - The type of enforcement for the policy.

* `exclude_resource_rules` - Exclude Resource Rules:

  * `api_groups` - List of API groups the resources belong to.

  * `api_versions` - List of API Versions the resources belong to.

  * `operations` - List of Operations the admission hook cares about.

  * `resources` - List of Resources this rule applies to.

* `failure_policy` - Failure policy to apply when the policy fails.

* `last_updated_at` - Date when the entity was last updated. The date is ISO 8601 and UTC.

* `last_updated_by` - The user the entity was last updated by.

* `match_conditions` - List of conditions that must be met for a request to be validated:

  * `expression` - Expression which will be evaluated by CEL.

  * `name` - Identifier for this match condition.

* `match_expressions` - List of label selector requirements that must be met for an object to be validated:

  * `key` - The label key that the selector applies to.

  * `operator` - A key's relationship to a set of values.

  * `values` - An array of string values.

* `match_labels` - Map of {key,value} pairs that must be met for an object to be validated.

* `match_policy` - Match policy.

* `name` - A human-friendly name used as an identifier for the policy instance.

* `org_id` - The id of the organization this entity belongs to.

* `project_criteria` - The project based criteria.

* `project_id` - The id of the project this entity belongs to.

* `resource_rules` - Resource Rules:

  * `api_groups` - List of API groups the resources belong to.

  * `api_versions` - List of API Versions the resources belong to.

  * `operations` - List of Operations the admission hook cares about.

  * `resources` - List of Resources this rule applies to.

* `validation_actions` - List of validation actions.

* `validations` - List of CEL expressions which are used to validate admission requests:

  * `expression` - Expression which will be evaluated by CEL.

  * `message` - Message displayed when validation fails.

  * `message_expression` - CEL expression that evaluates to the validation failure message that is returned when this rule fails.

  * `reason` - Machine-readable description of why this validation failed.
