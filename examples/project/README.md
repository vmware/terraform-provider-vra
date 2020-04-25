# project example

This is an example on how to create a project and assign a zone to it. 

## Getting Started

There are variables which need to be added to terraform.tfvars. The first are for connecting to the vRealize Automation (vRA) endpoint. There are names of cloud_account, region, zone, already setup in vRA.

* `url` - The URL for the vRealize Automation (vRA) endpoint
* `refresh_token` - The refresh token for the vRA account
* `insecure` - `false` for vRA Cloud and `true` for vRA on-prem
* `zone_name` - The compute placement zone within a region where machines can be placed
* `project_name` - Project name

To facilitate adding these variables, a sample tfvars file can be copied first:
```shell
cp terraform.tfvars.sample terraform.tfvars
```

A project can be created just with a project name. But this example shows how to create a project and assign zones to it. In order to achieve that, a cloud account must be setup in vRA and a zone must be created  within a region.

Follow these examples for setting up specific cloud accounts:

* Setup [cloud\_account\_aws](../cloud_account_aws/README.md)
* Setup [cloud\_account\_azure](../cloud_account_azure/README.md)
* Setup [cloud\_account\_gcp](../cloud_account_gcp/README.md)
* Setup [cloud\_account\_vsphere](../cloud_account_vsphere/README.md)

While the cloud account examples included setting up zones, here is an example
to setup a zone:

* Setup [zone](../zone/README.md)

Once the information is added to `terraform.tfvars`, the project can be created via:

```shell
terraform init
terraform apply
```
