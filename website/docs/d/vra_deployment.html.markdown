---
layout: "vra"
page_title: "VMware vRealize Automation: Data source vra_deployment"
description: A deployment data source.
---

# Data Source: vra\_deployment

This data source provides information about a deployment in vRA.

## Example Usages

This is an example of how to get a vRA deployment by its name.

```hcl
data "vra_deployment" "this" {
  name = var.deployment_name
}
```

This is an example of how to get a vRA cloud template by its id.

```hcl
data "vra_deployment" "this" {
  id = var.deployment_id
}
```


## Argument Reference

* `expand_last_request` - (Optional) Flag to indicate whether to expand last request on the deployment.

* `expand_project` - (Optional) Flag to indicate whether to expand project information.

* `expand_resources` - (Optional) Flag to indicate whether to expand resources in the deployment.

* `id` - (Optional) The id of the deployment. One of `id` or `name` must be provided.

* `name` - (Optional) The name of the deployment. One of `id` or `name` must be provided.


## Attribute Reference

* `blueprint_id` - The id of the cloud template used to request the deployment.

* `blueprint_version` - The version of the cloud template used to request the deployment.

* `catalog_item_id` - The id of the catalog item used to request the deployment.

* `catalog_item_version` - The version of the catalog item used to request the deployment.

* `created_at` - Date when the entity was created. The date is in ISO 6801 and UTC.

* `created_by` - The user the entity was created by.

* `description` - A human-friendly description.

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

* `inputs` - Inputs provided by the user while requesting / updating the deployment.

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
    
* `last_updated_at` - Date when the entity was last updated. The date is in ISO 6801 and UTC.

* `last_updated_by` - The user that last updated the deployment.

* `lease_expire_at` - Date when the deployment lease expire. The date is in ISO 6801 and UTC.

* `org_id` - The Id of the organization this deployment belongs to.

* `owner` - The user this deployment belongs to.

* `project` - The project this entity belongs to.

    * `description` - A human friendly description.
    
    * `id` - Id of the entity.
    
    * `name` - Name of the entity.
    
    * `version` - Version of the entity, if applicable.

* `project_id` - The id of the project this deployment belongs to.

* `resources` - Expanded resources for the deployment. Content of this property will not be maintained backward compatible.

    * `created_at` - Creation time (e.g. date format ‘2019-07-13T23:16:49.310Z’).
    
    * `depends_on` - A list of other resources this resource depends on.
    
    * `description` - A description of the resource.
    
    * `expense` - Expense incurred for this resource. 
    
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
