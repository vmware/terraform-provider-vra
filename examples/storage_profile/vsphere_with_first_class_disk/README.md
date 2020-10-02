# storage_profile example for vSphere with standard disk

This is an example on how to create a vSphere Storage Profile with first class disk. 

## Getting Started

There are variables which need to be added to terraform.tfvars. The first are for connecting to the vRealize Automation (vRA) endpoint. The others are names of cloud_account, and region that are already setup in vRA.

Following is the full list of variables used in this example.

* `url` - The URL for the vRealize Automation (vRA) endpoint
* `refresh_token` - The refresh token for the vRA account
* `cloud_account` - The name of the cloud account added in vRA
* `region` - The region within the cloud account
* `datastore_name` - The name of the vSphere fabric datastore
* `storage_policy_name` - The name of the vSphere fabric storage policy

To facilitate adding these variables, a sample tfvars file can be copied first:
```shell
cp terraform.tfvars.sample terraform.tfvars
```

Follow this example for setting up VMC and vSphere cloud accounts:

* Setup [cloud\_account\_vmc](../../cloud_account_vmc/README.md)
* Setup [cloud\_account\_vsphere](../../cloud_account_vsphere/README.md)

Once the information is added to `terraform.tfvars`, the Storage Profile can be created via:

```shell
terraform init
terraform apply
```
