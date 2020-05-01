# publish a catalog item from blueprint and request a deployment example

This is an example of how to publish a Cloud Assembly blueprint into Service Broker catalog and create a deployment from catalog in vRealize Automation(vRA).

Before requesting the deployment, a new project, a new blueprint are created. Besides creating the blueprint, it is versioned and released. Also a new blueprint content source and content sharing are created in order to publish the Cloud Assembly blueprint into Service Broker Catalog.

## Getting Started

There are variables which need to be added to terraform.tfvars.

* `url` - The URL for the vRealize Automation (vRA) endpoint
* `refresh_token` - The refresh token for the vRA account
* `insecure` - `false` for vRA Cloud and `true` for vRA on-prem
* `zone_name` - Existing Cloud Zone Name
* `project_name` - Project Name to create a new project
* `blueprint_name` - Blueprint name to create a new blueprint
* `catalog_source_name` - Catalog Source name to create a new content source for Cloud Assembly blueprints
* `deployment_name` - Deployment Name to request a new deployment

To facilitate adding these variables, a sample tfvars file can be copied first:
```shell
cp terraform.tfvars.sample terraform.tfvars
```

Once the information is added to `terraform.tfvars`, the cloud account can be brought up via:

```shell
terraform init
terraform apply
```
