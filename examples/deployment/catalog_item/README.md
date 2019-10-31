# deployment with a catalog item example

This is an example on how to request/create a deployment in vRealize Automation(vRA) using an existing catalog item.

## Getting Started

There are variables which need to be added to terraform.tfvars.

* `url` - The URL for the vRealize Automation (vRA) endpoint
* `refresh_token` - The refresh token for the vRA account
* `project_name` - Project Name
* `catalog_item_name` - Catalog Item Name
* `catalog_item_version` - Catalog Item Version
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
