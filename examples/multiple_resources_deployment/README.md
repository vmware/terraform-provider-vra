# Deployment with multiple resources example

This is an example on how to create multiple resources within the same deployment.
There are multiple ways of creating a deployment, for instance, through a blueprint, a catalog item or creating any resource like machine, network, load balancer or block device creates a deployment.
Associate the deployment id created through any of the above mentioned ways to the other resources that you want to create.

## Getting Started

There are variables which need to be added to terraform.tfvars. The first are for connecting to the vRealize Automation (vRA) endpoint. There are names of cloud_account, region, zone, project, already setup in vRA.

* `url` - The URL for the vRealize Automation (vRA) endpoint
* `refresh_token` - The refresh token for the vRA account
* `cloud_account` - The name of the cloud account added in vRA
* `region` - The region within in the cloud account
* `zone` - The compute placement zone within a region where machines can be placed
* `project` - The name of the project the current user belongs to
* `image_name` - The name of the fabric image corresponding to the cloud endpoint, such as ami-id for AWS.

To facilitate adding these variables, a sample tfvars file can be copied first:
```shell
cp terraform.tfvars.sample terraform.tfvars
```

To create a cloud-agnostic machine, a cloud account, zone, project, image and flavor profile must be setup. This is an example to create multiple resources within a single deployment. In this example, a block device (disk1) is created and its deployment id is passed into disk2, machine1 and machine2 

Follow these examples for setting up specific cloud accounts:

* Setup [cloud\_account\_aws](../cloud_account_aws/README.md)
* Setup [cloud\_account\_azure](../cloud_account_azure/README.md)
* Setup [cloud\_account\_gcp](../cloud_account_gcp/README.md)
* Setup [cloud\_account\_vsphere](../cloud_account_vsphere/README.md)

While the cloud account examples included setting up zones, here is an example
to setup a zone:

* Setup [zone](../zone/README.md)

To create a project, here is an example

* Setup [project](../project/README.md)

To create a deployment through blueprint or a catalog item, here is an example

* Setup [deployment](../deployment/README.md)

Once the information is added to `terraform.tfvars`, the image profile, flavor profile and machine can be created via:

```shell
terraform init
terraform apply
```
