# Content Source creation to pull a blueprint from a repository

This is an example on how to create a content source linking an already existing Integration with GitLab in vRealize Automation(vRA) to pull a blueprint from repository. This allows the use of Git(Lab|Hub) as a source for bluepints and ABX scripts.

This example pulls the blueprint from the blueprint01 folder in the repo at https://gitlab.com/vracontent/vra8_content_source_test

## Getting Started

There are variables which need to be added to terraform.tfvars.

* `url` - The URL for the vRealize Automation (vRA) endpoint
* `refresh_token` - The refresh token for the vRA account
* `insecure` - `true` for vRA on-prem and `false` for vRA Cloud
* `integration_id` - ID of the Git* integration to be used to connect to the repository/content
* `project_name` - Project Name

To facilitate adding these variables, a sample tfvars file can be copied first:
```shell
cp terraform.tfvars.sample terraform.tfvars
```

Once the information is added to `terraform.tfvars`, the cloud account can be brought up via:

```shell
terraform init
terraform apply
```
