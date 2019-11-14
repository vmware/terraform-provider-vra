# deployment with a blueprint id example

This is an example on how to crete a deployment in vRealize Automation(vRA) using an existing blueprint.

## Getting Started

There are variables which need to be added to terraform.tfvars.

* `url` - The URL for the vRealize Automation (vRA) endpoint
* `refresh_token` - The refresh token for the vRA account
* `project_name` - Project Name
* `blueprint_id` - Blueprint ID
* `blueprint_version` - Blueprint Version. Leave it empty for deployments with current draft of the blueprint.
* `deployment_name` - Deployment Name

To facilitate adding these variables, a sample tfvars file can be copied first:
```shell
cp terraform.tfvars.sample terraform.tfvars
```

Once the information is added to `terraform.tfvars`, the cloud account can be brought up via:

```shell
terraform init
terraform apply
```
