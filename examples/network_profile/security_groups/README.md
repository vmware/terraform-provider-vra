# network_profile example

This is an example on how to create a Network Profile with security groups so that firewall rules are added to all the machines provisioned with this Network Profile. 

## Getting Started

There are variables which need to be added to terraform.tfvars. The first are for connecting to the vRealize Automation (vRA) endpoint. The others are names of cloud_account, and region that already setup in vRA.

Following is the full list of variables used in this example.

* `url` - The URL for the vRealize Automation (vRA) endpoint
* `refresh_token` - The refresh token for the vRA account
* `cloud_account` - The name of the cloud account added in vRA
* `region` - The region within the cloud account
* `subnet_name` - The subnet to add to Network Profile
* `security_group_name` - The name of the security group to add to Network Profile

To facilitate adding these variables, a sample tfvars file can be copied first:
```shell
cp terraform.tfvars.sample terraform.tfvars
```

A Network Profile can be created just with a name and region id. But this example shows how to create a Network Profile and assign fabric networks (subnets) to it without isolation. In order to achieve that, a cloud account must be setup in vRA and the region must be identified within which the Network Profile to be created.

Follow these examples for setting up specific cloud accounts:

* Setup [cloud\_account\_aws](../../cloud_account_aws/README.md)
* Setup [cloud\_account\_azure](../../cloud_account_azure/README.md)
* Setup [cloud\_account\_gcp](../../cloud_account_gcp/README.md)
* Setup [cloud\_account\_vsphere](../../cloud_account_vsphere/README.md)

Once the information is added to `terraform.tfvars`, the Network Profile can be created via:

```shell
terraform init
terraform apply
```
