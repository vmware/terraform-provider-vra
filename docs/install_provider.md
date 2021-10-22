# Installing the Terraform Provider for VMware vRealize Automation

![Terraform](https://img.shields.io/badge/Terraform-0.13%2B-blue?style=for-the-badge&logo=terraform)

This document assumes the use of Terraform 0.13 or later.

## Automated Installation (Recommended)

The Terraform Provider for VMware vRealize Automation is a verified provider. Verified providers are owned and maintained by members of the HashiCorp Technology Partner Program. HashiCorp verifies the authenticity of the publisher and the providers are listed on the [Terraform Registry](https://registry.terraform.io) with a verified tier label. 

### Configure the Terraform Configuration Files

Providers listed on the Terraform Registry can be automatically downloaded when initializing a working directory with `terraform init`. The Terraform configuration block is used to configure some behaviors of Terraform itself, such as the Terraform version and the required providers and versions.

**Example**: A Terraform configuration block.

```hcl
terraform {
  required_providers {
    vra = {
      source  = "vmware/vra"
      version = ">= 0.4.0"
    }
  }
  required_version = ">= 0.13"
}
```
### Verify Terraform Initialization Using the Terraform Registry

To verify the initialization, navigate to the working directory for your Terraform configuration and run `terraform init`. You should see a message indicating that Terraform has been successfully initialized and downloaded the Terraform Provider for vRealize Automation from the Terraform Registry is installed.

**Example**: Initialize and Download the Provider.

```
$ ./terraform init

Initializing the backend...

Initializing provider plugins...
- Finding vmware/vra versions matching ">= 0.4.0"...
- Installing vmware/vra v0.4.0...
- Installed vmware/vra v0.4.0 (signed by a HashiCorp partner, key ID *************)

...

Terraform has been successfully initialized!
```

## Manual Installation

The [latest](https://github.com/vmware/terraform-provider-vra/releases/latest) release of the provider can be found on [the GitHub repository releases](https://github.com/vmware/terraform-provider-vra/releases). You can download the appropriate version of the provider for your operating system using a command line shell or a browser. 

This can be useful in environments that do not allow direct access to the Internet.

### Linux

The following examples use Bash on Linux (x64).

1. On an a Linux operating system with Internet access, download the plugin from GitHub using the shell. 

    ```bash
    RELEASE=0.4.0
    wget -q https://github.com/vmware/terraform-provider-vra/releases/download/v${RELEASE}/terraform-provider-vra_${RELEASE}_linux_amd64.zip
    ```

2. Extract the plugin.

    ```bash
    tar xvf terraform-provider-vra_${RELEASE}_linux_amd64.zip
    ```

3. Copy the extracted plugin to a target system and move to the Terraform plugins directory.

    >**Note**: The directory directory hierarchy that Terraform use to precisely determine the source of each provider it finds locally.<br/>
    > `$PLUGIN_DIRECTORY/$SOURCEHOSTNAME/$SOURCENAMESPACE/$NAME/$VERSION/$OS_$ARCH/`

    ```bash
    mv terraform-provider-vra_v${RELEASE} ~/.terraform.d/plugins/local/vmware/vra/${RELEASE}/linux_amd64
    ```

4. Verify the presence of the plugin in the Terraform plugins directory.

    ```bash
    cd ~/.terraform.d/plugins/local/vmware/vra/${RELEASE}/linux_amd64
    ls
    ```

### macOS

The following example uses Bash (default) on macOS (Intel).

1. On a macOS operating system with Internet access, install wget with [Homebrew](https://brew.sh).

    ```bash
    brew install wget
    ```

2. Download the plugin from GitHub using the shell. 

    ```bash
    RELEASE=0.4.0
    wget -q https://github.com/vmware/terraform-provider-vra/releases/download/v${RELEASE}/terraform-provider-vra_${RELEASE}_darwin_amd64.zip
    ```

3. Extract the plugin.

    ```bash
    tar xvf terraform-provider-vra_${RELEASE}_darwin_amd64.zip
    ```

5. Copy the extracted plugin to a target system and move to the Terraform plugins directory.

    >**Note**: The directory directory hierarchy that Terraform use to precisely determine the source of each provider it finds locally.<br/>
    > `$PLUGIN_DIRECTORY/$SOURCEHOSTNAME/$SOURCENAMESPACE/$NAME/$VERSION/$OS_$ARCH/`

    ```bash
    mv terraform-provider-vra_v${RELEASE} ~/.terraform.d/plugins/local/vmware/vra/${RELEASE}/darwin_amd64
    ```

6. Verify the presence of the plugin in the Terraform plugins directory.

    ```bash
    cd ~/.terraform.d/plugins/local/vmware/vra/${RELEASE}/darwin_amd64
    ls
    ```

### Windows

The following examples use PowerShell on Windows (x64).

1. On a Windows operating system with Internet access, download the plugin using the PowerShell. 

    ```powershell
    $RELEASE="0.4.0"
    Invoke-WebRequest https://github.com/vmware/terraform-provider-vra/releases/download/v${RELEASE}/terraform-provider-vra_${RELEASE}_windows_amd64.zip -outfile terraform-provider-vra_${RELEASE}_windows_amd64.zip
    ```

2. Extract the plugin.

    ```powershell
    Expand-Archive terraform-provider-vra_${RELEASE}_windows_amd64.zip

    cd terraform-provider-vra_${RELEASE}_windows_amd64
    ```

4. Copy the extracted plugin to a target system and move to the Terraform plugins directory.

    >**Note**: The directory directory hierarchy that Terraform use to precisely determine the source of each provider it finds locally.<br/>
    > `$PLUGIN_DIRECTORY/$SOURCEHOSTNAME/$SOURCENAMESPACE/$NAME/$VERSION/$OS_$ARCH/`

    ```powershell
    New-Item $ENV:APPDATA\terraform.d\plugins\local\vmware\vra\${RELEASE}\ -Name "windows_amd64" -ItemType "directory"

    Move-Item terraform-provider-vra_v${RELEASE}.exe $ENV:APPDATA\terraform.d\plugins\local\vmware\vra\${RELEASE}\windows_amd64\terraform-provider-vra_v${RELEASE}.exe
    ```

5. Verify the presence of the plugin in the Terraform plugins directory.

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
      version = ">= 0.4.0"
    }
  }
  required_version = ">= 1.0.0"
}
```

### Verify the Terraform Initialization of a Manually Installed Provider

To verify the initialization, navigate to the working directory for your Terraform configuration and run `terraform init`. You should see a message indicating that Terraform has been successfully initialized and the installed version of the Terraform Provider for vRealize Automation.

**Example**: Initialize and Use a Manually Installed Provider

```
$ ./terraform init

Initializing the backend...

Initializing provider plugins...
- Finding local/vmware/vra versions matching ">= 0.4.0"...
- Installing local/vmware/vra v0.4.0...
- Installed local/vmware/vra v0.4.0 (unauthenticated)
...

Terraform has been successfully initialized!
```

## Get the Provider Version
To find the provider version, navigate to the working directory of your Terraform configuration and run `terraform version`. You should see a message indicating the provider version.

**Example**: Terraform Provider Version from the Terraform Registry 

```
$ ./terraform version
Terraform v1.0.0
on linux_amd64
+ provider registry.terraform.io/vmware/vra v0.4.0
```
**Example**: Terraform Provider Version for a Manually Installed Provider

```
$ ./terraform version
Terraform v1.0.0
on linux_amd64
+ provider local/vmware/vra v0.4.0
```
