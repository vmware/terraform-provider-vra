# project example

This is an example on how to create a project and assign a zone to it. 

## Getting Started

There are variables which need to be added to terraform.tfvars. The first are for connecting to the vRealize Automation (vRA) endpoint, there are credentials for connecting to the vCenter instance, and the name of the data collector already setup in vRA.

* `url` - The URL for the vRealize Automation (vRA) endpoint
* `refresh_token` - The refresh token for the vRA account
* `cloud_account` - The name of the cloud account added in vRA
* `region` - The region within in the cloud account
* `zone` - The compute placement zone within a region where machines can be placed

To facilitate adding these variables, a sample tfvars file can be copied first:
```shell
cp terraform.tfvars.sample terraform.tfvars
```

A project can be created just with a project name. But this example shows how to create a project and assign zones to it. In order to achieve that, a cloud account must be setup in vRA and a zone must be created  within a region.

Follow these examples for setting up specific cloud accounts:

* Setup [cloud\_account\_aws](examples/cloud_account_aws/README.md)
* Setup [cloud\_account\_azure](examples/cloud_account_azure/README.md)
* Setup [cloud\_account\_vsphere](examples/cloud_account_vsphere/README.md)

While the cloud account examples included setting up zones, here is an example
to setup a zone:

* Setup [zone](examples/zone/README.md)

Once the information is added to `terraform.tfvars`, the zones can be created via:

```shell
terraform init
terraform apply
```
