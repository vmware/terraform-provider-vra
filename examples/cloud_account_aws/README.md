# cloud\_account\_aws example

This is an example on how to setup an AWS cloud account.

## Getting Started

There are variables which need to be added to terraform.tfvars. The first are for connecting to the vRealize Automation (vRA) endpoint, there are credentials for connecting to the AWS cloud account.

* `url` - The URL for the vRealize Automation (vRA) endpoint
* `refresh_token` - The refresh token for the vRA account
* `access_key` - AWS access key ID
* `secret_key` - AWS secret access key

To facilitate adding these variables, a sample tfvars file can be copied first:
```shell
cp terraform.tfvars.sample terraform.tfvars
```

Once the information is added to `terraform.tfvars`, the cloud account can be brought up via:

```shell
terraform init
terraform apply
```
