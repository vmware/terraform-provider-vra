# deployment with a blueprint id example

This is an example on how to crete a blueprint in vRealize Automation(vRA).

## Getting Started

There are variables which need to be added to terraform.tfvars.

* `url` - The URL for the vRealize Automation (vRA) endpoint
* `refresh_token` - The refresh token for the vRA account
* `insecure` - `true` for vRA on-prem and `false` for vRA Cloud
* `project_name` - Project Name
* `blueprint_name` - Blueprint Name

To facilitate adding these variables, a sample tfvars file can be copied first:
```shell
cp terraform.tfvars.sample terraform.tfvars
```

Once the information is added to `terraform.tfvars`, the cloud account can be brought up via:

```shell
terraform init
terraform apply
```
