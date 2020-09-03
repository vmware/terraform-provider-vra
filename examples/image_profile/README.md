# image_profile example

This is an example on how to create an image profile with an image id in vRealie Automation.
Image profile represents a structure that holds a list of image mappings defined for the particular region.
Flavor profile represents a structure that holds flavor mappings defined for the corresponding cloud end-point region.

## Getting Started

There are variables which need to be added to terraform.tfvars. The first are for connecting to the vRealize Automation (vRA) endpoint. There are names of cloud_account, region, zone, project, already setup in vRA.

* `url` - The URL for the vRealize Automation (vRA) endpoint
* `refresh_token` - The refresh token (API token) for the vRA user account
* `cloud_account` - The name of the cloud account added in vRA
* `region` - The region within in the cloud account. For vSphere, it is the externalRegionId such as `Datacenter:datacenter-2` and for AWS, it is region name such as `us-east-1`, etc. 
* `image_name1` - The name of the fabric image corresponding to the cloud endpoint, such as ami-id for AWS, template name for vSphere, etc.
* `image_name2` - The name of the fabric image corresponding to the cloud endpoint, such as ami-id for AWS, template name for vSphere, etc.

To facilitate adding these variables, a sample tfvars file can be copied first:
```shell
cp terraform.tfvars.sample terraform.tfvars
```

Follow these examples for setting up specific cloud accounts:

* Setup [cloud\_account\_aws](../cloud_account_aws/README.md)
* Setup [cloud\_account\_azure](../cloud_account_azure/README.md)
* Setup [cloud\_account\_gcp](../cloud_account_gcp/README.md)
* Setup [cloud\_account\_vmc](../cloud_account_vmc/README.md)
* Setup [cloud\_account\_vsphere](../cloud_account_vsphere/README.md)

Once the information is added to `terraform.tfvars`, the image profile can be created via:

```shell
terraform init
terraform apply
```
