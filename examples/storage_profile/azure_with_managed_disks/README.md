# storage_profile example for Azure with managed disk storage type

This is an example on how to create an Azure Storage Profile with managed disk storage type. 

## Getting Started

There are variables which need to be added to terraform.tfvars. The first are for connecting to the vRealize Automation (vRA) endpoint. The others are names of cloud_account, and region that are already setup in vRA.

Following is the full list of variables used in this example.

* `url` - The URL for the vRealize Automation (vRA) endpoint
* `refresh_token` - The refresh token for the vRA account. Alternatively, use the VRA_REFRESH_TOKEN 
* `cloud_account` - The name of the Azure cloud account added in vRA
* `region` - The region within the Azure cloud account

To facilitate adding these variables, a sample tfvars file can be copied first:
```shell
cp terraform.tfvars.sample terraform.tfvars
```

Follow this examples for setting up Azure cloud account:

* Setup [cloud\_account\_azure](../../cloud_account_azure/README.md)

Once the information is added to `terraform.tfvars`, the Storage Profile can be created via:

```shell
terraform init
terraform apply
```
