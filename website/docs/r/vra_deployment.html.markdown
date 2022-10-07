---
layout: "vra"
page_title: "VMware vRealize Automation: Resource vra_deployment"
description: A resource that can be used to create a vRealize Automation deployment.
---

# Resource: vra\_deployment

This resource provides a way to create a deployment in vRealize Automation(vRA) by either using a blueprint, or catalog item, or an inline blueprint.

## Example Usages

This is an example of how to create a deployment using a catalog item.

```hcl
resource "vra_deployment" "this" {
  name        = var.deployment_name
  description = "Deployment description"

  catalog_item_id      = var.catalog_item_id
  catalog_item_version = var.catalog_item_version
  project_id           = var.project_id

  inputs = {
    flavor = "small"
    image  = "centos"
    count  = 1
    flag   = false
    number = 10.0
    arrayProp = jsonencode(["foo", "bar", "where", "waldo"])
    objectProp = jsonencode({ "key1": "value1", "key2": [1, 2, 3, 4] })
  }

  timeouts {
    create = "30m"
    delete = "30m"
    update = "30m"
  }
}
```

This is an example of how to create a deployment using a cloud template.

```hcl
resource "vra_deployment" "this" {
  name        = var.deployment_name
  description = "Deployment description"

  blueprint_id      = var.blueprint_id
  blueprint_version = var.blueprint_version
  project_id        = var.project_id

  inputs = {
    flavor = "small"
    image  = "centos"
    count  = 1
    flag   = true
    arrayProp = jsonencode(["foo", "bar", "baz"])
    objectProp = jsonencode({ "key": "value", "key2": [1, 2, 3] })
  }

  timeouts {
    create = "30m"
    delete = "30m"
    update = "30m"
  }
}
```

This is an example of how to create a deployment without any resources so that it may be attached to other IaaS resources like `vra_machine`, `vra_network`, etc.

```hcl
resource "vra_deployment" "this" {
  name        = var.deployment_name
  description = "Deployment description"

  project_id = var.project_id
}
```

## Argument Reference

* `blueprint_id` - (Optional) The id of the cloud template to be used to request the deployment. Conflicts with `blueprint_content` and `catalog_item_id`.

* `blueprint_version` - (Optional) The version of the cloud template to be used to request the deployment. Used only when `blueprint_id` is provided.

* `blueprint_content` - (Optional) The content of the the cloud template to be used to request the deployment. Conflicts with `blueprint_id` and `catalog_item_id`.

* `catalog_item_id` - (Optional) The id of the catalog item to be used to request the deployment. Conflicts with `blueprint_id` and `blueprint_content`.

* `catalog_item_version` - (Optional) The version of the catalog item to be used to request the deployment. Used only when `catalog_item_id` is provided.

* `description` - (Optional) A human-friendly description.

* `expand_project` - (Optional) Flag to indicate whether to expand project information.

* `inputs` - (Optional) Inputs provided by the user. For inputs including those with default values, refer to `inputs_including_defaults`.

* `name` - (Required) The name of the deployment.

* `owner` - (Optional) The user this deployment belongs to. At create, the owner is ignored but is used to update during next apply.

* `project_id` - (Required) The id of the project this deployment belongs to.

* `reason` - (Optional) Reason for requesting/updating a blueprint.

## Attribute Reference

* `created_at` - Date when the entity was created. The date is in ISO 6801 and UTC.

* `created_by` - The user the entity was created by.

* `expense` - Expense incurred for the deployment. 

    * `additional_expense` - Additional expense incurred for the resource.
    
    * `code` - Expense sync message code if any.
    
    * `compute_expense` - Compute expense of the entity.
    
    * `last_update_time` - Last expense sync time.
    
    * `message` - Expense sync message if any.
    
    * `network_expense` - Network expense of the entity.
    
    * `storage_expense` - Storage expense of the entity.
    
    * `total_expense` - Total expense of the entity.
    
    * `unit` - Monetary unit.
    
* `id` - The id of the deployment.

* `inputs_including_defaults` - All the inputs applied during last create/update operation, including those with default values. For the list of inputs provided by the user in the configuration, refer to `inputs`.

* `last_request` - Represents deployment requests.

    * `action_id` - Identifier of the requested action.
    
    * `approved_at` - Time at which the request was approved.
    
    * `blueprint_id` - Identifier of the requested blueprint in the form ‘UUID:version’.
    
    * `cancelable` - Indicates whether request can be canceled or not. 
    
    * `catalog_item_id` - Identifier of the requested catalog item in the form ‘UUID:version’.
    
    * `completed_at` - Time at which the request completed.
    
    * `completed_tasks` - The number of tasks completed while fulfilling this request.
    
    * `created_at` - Creation time (e.g. date format ‘2019-07-13T23:16:49.310Z’).
    
    * `details` - Longer user-friendly details of the request.
 
    * `dismissed` - Indicates whether request is in dismissed state.
     
    * `id` - Request identifier.
 
    * `initialized_at` - Time at which the request was initialized.

    * `inputs` - List of request inputs.
    
    * `name` - Short user-friendly label of the request (e.g. ‘shuting down myVM’).
    
    * `outputs` - Request outputs.
    
    * `requested_by` - The user that initiated the request.
     
    * `resource_name` - Optional resource name to which the request applies to.
 
    * `status` - Request overall execution status. Supported values: `CREATED`, `PENDING`, `INITIALIZATION`, `CHECKING_APPROVAL`, `APPROVAL_PENDING`, `INPROGRESS`, `COMPLETION`, `APPROVAL_REJECTED`, `ABORTED`, `SUCCESSFUL`, `FAILED`.

    * `total_tasks` -The total number of tasks need to be completed to fulfil this request.

    * `updated_at` - Last update time (e.g. date format ‘2019-07-13T23:16:49.310Z’).
    
* `last_updated_at` - TDate when the entity was last updated. The date is in ISO 6801 and UTC.

* `last_updated_by` - The user that last updated the deployment.

* `lease_expire_at` - Date when the deployment lease expire. The date is in ISO 6801 and UTC.

* `org_id` - The Id of the organization this deployment belongs to.

* `project` - The project this entity belongs to.

    * `description` - A human friendly description.
    
    * `id` - Id of the entity.
    
    * `name` - Name of the entity.
    
    * `version` - Version of the entity, if applicable.

* `resources` - Expanded resources for the deployment. Content of this property will not be maintained backward compatible.

    * `created_at` - Creation time (e.g. date format ‘2019-07-13T23:16:49.310Z’).
    
    * `depends_on` - A list of other resources this resource depends on.
    
    * `description` - A description of the resource.
    
    `expense` - Expense incurred for this resource. 
    
        * `additional_expense` - Additional expense incurred for the resource.
        
        * `code` - Expense sync message code if any.
        
        * `compute_expense` - Compute expense of the entity.
        
        * `last_update_time` - Last expense sync time.
        
        * `message` - Expense sync message if any.
        
        * `network_expense` - Network expense of the entity.
        
        * `storage_expense` - Storage expense of the entity.
        
        * `total_expense` - Total expense of the entity.
        
        * `unit` - Monetary unit.
    
    * `id` - Unique identifier of the resource.
    
    * `name` - Name of the resource.
    
    * `properties_json` - List of properties in the encoded JSON string format. 
    
    * `state` - The current state of the resource. Supported values are `PARTIAL`, `TAINTED`, `OK.`
    
    * `sync_status` - The current sync status. Supported values are `SUCCESS`, `MISSING`, `STALE`.
    
    * `type` - Type of the resource.

* `status` - The status of the deployment with respect to its life cycle operations.

## Import

Deployment can be imported using the id, e.g.

`$ terraform import vra_deployment.this 05956583-6488-4e7d-84c9-92a7b7219a15`
