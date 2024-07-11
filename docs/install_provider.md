# Installing the Terraform Provider for VMware Aria Automation

This document assumes the use of Terraform 0.13 or later.

## Automated Installation (Recommended)

The Terraform Provider for VMware Aria Automation is a partner provider. Partner providers are owned and maintained by members of the HashiCorp Technology Partner Program. HashiCorp verifies the authenticity of the publisher and the providers are listed on the [Terraform Registry](https://registry.terraform.io) with a `Partner` label.

### Configure the Terraform Configuration Files

Providers listed on the Terraform Registry can be automatically downloaded when initializing a working directory with `terraform init`. The Terraform configuration block is used to configure some behaviors of Terraform itself, such as the Terraform version and the required providers and versions.

**Example**: A Terraform configuration block.

```hcl
terraform {
  required_providers {
    vra = {
      source  = "vmware/vra"
    }
  }
  required_version = ">= 0.13"
}
```

You can use `version` locking and operators to require specific versions of the provider.

**Example**: A Terraform configuration block with the provider versions.

```hcl
terraform {
  required_providers {
    vra = {
      source  = "vmware/vra"
      version = ">= x.y.z"
    }
  }
  required_version = ">= 0.13"
}
```

To specify a particular provider version when installing released providers, see the Terrraform documentation [on provider versioning](https://www.terraform.io/docs/configuration/providers.html#version-provider-versions).

### Verify Terraform Initialization Using the Terraform Registry

To verify the initialization, navigate to the working directory for your Terraform configuration and run `terraform init`. You should see a message indicating that Terraform has been successfully initialized and has installed the provider from the Terraform Registry.

**Example**: Initialize and Download the Provider.

```shell
$ terraform init

Initializing the backend...

Initializing provider plugins...
- Finding vmware/vra versions matching ">= x.y.z" ...
- Installing vmware/vra x.y.z ...
- Installed vmware/vra x.y.z (signed by a HashiCorp partner, key ID *************)

...

Terraform has been successfully initialized!
```

## Manual Installation

The [latest](https://github.com/vmware/terraform-provider-vra/releases/latest) release of the provider can be found on [the GitHub repository releases](https://github.com/vmware/terraform-provider-vra/releases). You can download the appropriate version of the provider for your operating system using a command line shell or a browser.

This can be useful in environments that do not allow direct access to the Internet.

### Linux

The following examples use Bash on Linux (x64).

1. On an a Linux operating system with Internet access, download the plugin from GitHub using the shell.

    ```shell
    RELEASE=x.y.z
    wget -q https://github.com/vmware/terraform-provider-vra/releases/download/v${RELEASE}/terraform-provider-vra_${RELEASE}_linux_amd64.zip
    ```

2. Extract the plugin.

    ```shell
    unzip terraform-provider-vra_${RELEASE}_linux_amd64.zip
    ```

3. Create a directory for the provider.

    >**Note**: The directory hierarchy that Terraform use to precisely determine the source of each provider it finds locally.<br/>
    > `$PLUGIN_DIRECTORY/$SOURCEHOSTNAME/$SOURCENAMESPACE/$NAME/$VERSION/$OS_$ARCH/`

    ```shell
    mkdir -p ~/.terraform.d/plugins/local/vmware/vra/${RELEASE}/linux_amd64
    ```

4. Copy the extracted plugin to a target system and move to the Terraform plugins directory.

    ```shell
    mv terraform-provider-vra_v${RELEASE} ~/.terraform.d/plugins/local/vmware/vra/${RELEASE}/linux_amd64
    ```

5. Verify the presence of the plugin in the Terraform plugins directory.

    ```shell
    cd ~/.terraform.d/plugins/local/vmware/vra/${RELEASE}/linux_amd64
    ls
    ```

### macOS

The following example uses Bash (default) on macOS (Intel).

1. On a macOS operating system with Internet access, install wget with [Homebrew](https://brew.sh).

    ```shell
    brew install wget
    ```

2. Download the plugin from GitHub using the shell.

    ```shell
    RELEASE=x.y.x
    wget -q https://github.com/vmware/terraform-provider-vra/releases/download/v${RELEASE}/terraform-provider-vra_${RELEASE}_darwin_amd64.zip
    ```

3. Extract the plugin.

    ```shell
    tar xvf terraform-provider-vra_${RELEASE}_darwin_amd64.zip
    ```

4. Create a directory for the provider.

    >**Note**: The directory hierarchy that Terraform use to precisely determine the source of each provider it finds locally.<br/>
    > `$PLUGIN_DIRECTORY/$SOURCEHOSTNAME/$SOURCENAMESPACE/$NAME/$VERSION/$OS_$ARCH/`

    ```shell
    mkdir -p ~/.terraform.d/plugins/local/vmware/vra/${RELEASE}/darwin_amd64
    ```

5. Copy the extracted plugin to a target system and move to the Terraform plugins directory.

    ```shell
    mv terraform-provider-vra_v${RELEASE} ~/.terraform.d/plugins/local/vmware/vra/${RELEASE}/darwin_amd64
    ```

6. Verify the presence of the plugin in the Terraform plugins directory.

    ```shell
    cd ~/.terraform.d/plugins/local/vmware/vra/${RELEASE}/darwin_amd64
    ls
    ```

### Windows

The following examples use PowerShell on Windows (x64).

1. On a Windows operating system with Internet access, download the plugin using the PowerShell.

    ```powershell
    $RELEASE="x.y.z"
    Invoke-WebRequest https://github.com/vmware/terraform-provider-vra/releases/download/v${RELEASE}/terraform-provider-vra_${RELEASE}_windows_amd64.zip -outfile terraform-provider-vra_${RELEASE}_windows_amd64.zip
    ```

2. Extract the plugin.

    ```powershell
    Expand-Archive terraform-provider-vra_${RELEASE}_windows_amd64.zip

    cd terraform-provider-vra_${RELEASE}_windows_amd64
    ```

3. Copy the extracted plugin to a target system and move to the Terraform plugins directory.

    >**Note**: The directory directory hierarchy that Terraform use to precisely determine the source of each provider it finds locally.<br/>
    > `$PLUGIN_DIRECTORY/$SOURCEHOSTNAME/$SOURCENAMESPACE/$NAME/$VERSION/$OS_$ARCH/`

    ```powershell
    New-Item $ENV:APPDATA\terraform.d\plugins\local\vmware\vra\${RELEASE}\ -Name "windows_amd64" -ItemType "directory"

    Move-Item terraform-provider-vra_v${RELEASE}.exe $ENV:APPDATA\terraform.d\plugins\local\vmware\vra\${RELEASE}\windows_amd64\terraform-provider-vra_v${RELEASE}.exe
    ```

4. Verify the presence of the plugin in the Terraform plugins directory.

    ```powershell
    cd $ENV:APPDATA\terraform.d\plugins\local\vmware\vra\${RELEASE}\windows_amd64
    dir
    ```

### Configure the Terraform Configuration Files

A working directory can be initialized with providers that are installed locally on a system by using `terraform init`. The Terraform configuration block is used to configure some behaviors of Terraform itself, such as the Terraform version and the required providers source and version.

**Example**: A Terraform configuration block.

```hcl
terraform {
  required_providers {
    vra = {
      source  = "local/vmware/vra"
      version = ">= x.y.z"
    }
  }
  required_version = ">= 0.13"
}
```

### Verify the Terraform Initialization of a Manually Installed Provider

To verify the initialization, navigate to the working directory for your Terraform configuration and run `terraform init`. You should see a message indicating that Terraform has been successfully initialized and the installed version of the Terraform Provider for VMware Aria Automation.

**Example**: Initialize and Use a Manually Installed Provider

```shell
$ terraform init

Initializing the backend...

Initializing provider plugins...
- Finding local/vmware/vra versions matching ">= x.y.x" ...
- Installing local/vmware/vra x.y.x ...
- Installed local/vmware/vra x.y.x (unauthenticated)
...

Terraform has been successfully initialized!
```

## Get the Provider Version

To find the provider version, navigate to the working directory of your Terraform configuration and run `terraform version`. You should see a message indicating the provider version.

**Example**: Terraform Provider Version from the Terraform Registry

```shell
$ terraform version
Terraform x.y.z
on linux_amd64
+ provider registry.terraform.io/vmware/vra x.y.z
```

**Example**: Terraform Provider Version for a Manually Installed Provider

```shell
$ terraform version
Terraform x.y.z
on linux_amd64

+ provider local/vmware/vra x.y.z

```
