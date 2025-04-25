provider "vra" {
  url           = var.url
  refresh_token = var.refresh_token
  insecure      = var.insecure
}

// Approval Policy
resource "vra_policy_approval" "policy_approval" {
  name             = "terraform-approval-policy"
  description      = "Approval Policy [terraform-approval-policy] created by Terraform"
  enforcement_type = "HARD"
  project_criteria = [
    {
      key      = "project.name"
      operator = "eq"
      value    = "default-peoject"
    }
  ]

  actions = [
    "Deployment.ChangeLease",
  ]
  approval_level = 1
  approval_mode  = "ANY_OF"
  approval_type  = "ROLE"
  approvers = [
    // "USER:admin",
    // "GROUP:vraadamadmins@",
    "ROLE:PROJECT_ADMINISTRATORS"
  ]
  auto_approval_decision = "APPROVE"
  auto_approval_expiry   = 30
}

data "vra_policy_approval" "policy_approval" {
  depends_on = [
    vra_policy_approval.policy_approval
  ]

  search = "terraform-approval-policy"
}

output "policy_approval" {
  value = data.vra_policy_approval.policy_approval
}

// Day2 Action Policy
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
    //"ROLE:administrator"
  ]
}

data "vra_policy_day2_action" "policy_day2_action" {
  depends_on = [
    vra_policy_day2_action.policy_day2_action
  ]

  search = "terraform-day2-action-policy"
}

output "policy_day2_action" {
  value = data.vra_policy_day2_action.policy_day2_action
}

// IaaS Resource Policy
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

data "vra_policy_iaas_resource" "policy_iaas_resource" {
  depends_on = [
    vra_policy_iaas_resource.policy_iaas_resource
  ]

  search = "terraform-iaas-resource-policy"
}

output "policy_iaas_resource" {
  value = data.vra_policy_iaas_resource.policy_iaas_resource
}

// Lease Policy
resource "vra_policy_lease" "policy_lease" {
  name             = "terraform-lease-policy"
  description      = "Lease Policy [terraform-lease-policy] created by Terraform"
  enforcement_type = "HARD"

  criteria = [
    {
      "and" : jsonencode([
        {
          "key" : "ownedBy",
          "operator" : "eq",
          "value" : "admin"
        },
        {
          "or" : [
            {
              "key" : "createdBy",
              "operator" : "eq",
              "value" : "bob"
            },
            {
              "key" : "createdBy",
              "operator" : "eq",
              "value" : "jeff"
            }
          ]
        }
      ])
    }
  ]

  lease_grace          = 15
  lease_term_max       = 30
  lease_total_term_max = 100
}

data "vra_policy_lease" "policy_lease" {
  depends_on = [
    vra_policy_lease.policy_lease
  ]

  search = "terraform-lease-policy"
}

output "policy_lease" {
  value = data.vra_policy_lease.policy_lease
}
