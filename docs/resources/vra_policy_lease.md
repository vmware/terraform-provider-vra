---
page_title: "VMware Aria Automation: vra_policy_lease"
description: A resource for Lease policies.
---

# Resource: vra_policy_lease

Creates a Lease policy resource to automate the expiration and destruction of deployed catalog items.

## Example Usages

The following example shows how to create a lease policy resource:

```hcl
resource "vra_policy_lease" "policy_lease" {
  name             = "terraform-lease-policy"
  description      = "Lease Policy [terraform-lease-policy] created by Terraform"
  enforcement_type = "HARD"

  lease_grace          = 15
  lease_term_max       = 30
  lease_total_term_max = 100
}
```

## Argument Reference

Create your resource with the following arguments:

* `criteria` - (Optional) The policy criteria.

* `description` - (Optional) A human-friendly description for the policy instance.

* `enforcement_type` - (Required) The type of enforcement for the policy. Supported values: `HARD`, `SOFT`.

* `lease_grace` - (Optional) The duration in days that an expired object should be held before it is deleted. Valid range: `0` - `127`.

* `lease_term_max` - (Required) The maximum duration in days between creation (or renewal) and expiration. Valid range: `1` - `32767`.

* `lease_total_term_max` - (Required) The maximum duration in days between creation and expiration. Unaffected by renewal. Valid range: `1` - `32767`.

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

To import an existing Lease policy, use the `id` as in the following example:

`$ terraform import vra_policy_lease.policy_lease "39616df1-f42c-4ef1-a8e1-1a2abfec1fd6"`
