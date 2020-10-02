# cloud\_account\_azure example

This is an example on how to setup an Azure cloud account.

## Getting Started

There are variables which need to be added to terraform.tfvars. The first are for connecting to the vRealize Automation (vRA) endpoint, there are credentials for connecting to the Azure cloud account.

* `url` - The URL for the vRealize Automation (vRA) endpoint
* `refresh_token` - The refresh token of the vRA account
* `subscription_id  ` - Azure subscription id
* `tenant_id` - Azure tenant id
* `application_id` - Azure client application id
* `application_key` - Azure client application secret key

To facilitate adding these variables, a sample tfvars file can be copied first:
```shell
cp terraform.tfvars.sample terraform.tfvars
```

Once the information is added to `terraform.tfvars`, the cloud account can be brought up via:

```shell
terraform init
terraform apply
```
