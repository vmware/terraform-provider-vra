# network_ip_range example

This is an example on how to create a Network IP Range associated with a network in vRA that can used to let vRA assign IP addresses if not using an external IPAM solution like InfoBlox or Bluecat

## Getting Started

There are variables which need to be added to terraform.tfvars. The first two are for connecting to the vRealize Automation (vRA) endpoint. The next two are the names of a cloud_account, and subnet that already setup in vRA. The last three are the actual values needed to create the ip range. 

Following is the full list of variables used in this example.

* `url` - The URL for the vRealize Automation (vRA) endpoint
* `refresh_token` - The refresh token for the vRA account
* `cloud_account` - The name of the cloud account in vRA that the network is attached to
* `subnet_name` - The name of the subnet to associate the ip range with. *Must have CIDR info and gateway populated*
* `start_ip` - First IP of the range. Can not include the default gateway.
* `end_ip` - Last IP of the range
* `ip_version` - IP networking version either IPv4 or IPv6

To facilitate adding these variables, a sample tfvars file can be copied first:
```shell
cp terraform.tfvars.sample terraform.tfvars
```

A Network IP Range requires that the network it's being associated with have CIDR information and the default gateway assigned. 

Follow these examples for setting up specific cloud accounts:

* Setup [cloud\_account\_aws](../../cloud_account_aws/README.md)
* Setup [cloud\_account\_azure](../../cloud_account_azure/README.md)
* Setup [cloud\_account\_gcp](../../cloud_account_gcp/README.md)
* Setup [cloud\_account\_vsphere](../../cloud_account_vsphere/README.md)

Once the information is added to `terraform.tfvars`, the Network IP Range can be created via:

```shell
terraform init
terraform apply
```
