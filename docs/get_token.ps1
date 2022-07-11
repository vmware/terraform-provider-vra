<#
    .NOTES
    THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE
    WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
    COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR
    OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.    

    .SYNOPSIS
        Generates and returns a `refresh_token` from vRealize Automation Cloud or vRealize Automation for use by the Terraform provider.

    .DESCRIPTION
        The Request-vRARefreshToken function connects to the specified vRealize Automation endpoint and obtains a refresh token that is needed the Terraform provider.
            terraform {
                required_providers {
                    vra = {
                        source  = "vmware/vra"
                        version = ">= x.y.z"
                }
            }
                required_version = ">= 0.13"
            }
            provider "vra" {
                url           = "https://api.mgmt.cloud.vmware.com"
                refresh_token = "mx7w9**********************zB3UC"
                insecure      = false
            }

    .PARAMETER FQDN
        (string) FQDN of the vRA Instance
    
    .PARAMETER Username
        (string) Username used to connect to the vRA instance

    .PARAMETER Password
        (string) Password used to authenticate
    
    .PARAMETER Domain
        (string) Authentication domain configured
    
    .PARAMETER SkipCertValidation
        (switch) Skip certificate validation

    .EXAMPLE
        get_token.ps1
        #Fully interactive
        Enter the FQDN for the vRealize Automation services: cloud.rainpole.io
        Enter the username to authenticate with vRealize Automation: john.doe
        Enter the password to authenticate with vRealize Automation: ********
        Enter the domain or press enter to skip: example.com
        Successfully connected to endpoint for vRealize Automation services: cloud.rainpole.io
        Generating token...

        ---------Refresh Token---------
        N5WD************************kfgm
        -------------------------------

        Saving environmental variables...

        VRA_URL = https://cloud.rainpole.io
        VRA_REFRESH_TOKEN = N5WD************************kfgm

    .EXAMPLE
        get_token.ps1 -FQDN cloud.rainpole.io -Username john.doe -Password ****** -Domain example.com
        #All parameters on CLI

        Successfully connected to endpoint for vRealize Automation services: cloud.rainpole.io
        Generating token...

        ---------Refresh Token---------
        N5WD************************kfgm
        -------------------------------

        Saving environmental variables...

        VRA_URL = https://cloud.rainpole.io
        VRA_REFRESH_TOKEN = N5WD************************kfgm

    .EXAMPLE
        get_token.ps1 -FQDN cloud.rainpole.io -Username john.doe -Domain example.com
        # Enter only password interactively to avoid it entering console history
        Enter the password to authenticate with vRealize Automation: *****************
        Successfully connected to endpoint for vRealize Automation services: cloud.rainpole.io
        Generating token...

        ---------Refresh Token---------
        N5WD************************kfgm
        -------------------------------

        Saving environmental variables...

        VRA_URL = https://cloud.rainpole.io
        VRA_REFRESH_TOKEN = N5WD************************kfgm
    
    .INPUTS
    None

    .OUTPUTS
    Text

    .LINK
    https://github.com/vmware/terraform-provider-vra

#>
[CmdletBinding()]
Param (
    [Parameter ()]
        [ValidateNotNullOrEmpty()]
        [Alias('Server')] 
        [String]$FQDN = (Read-Host "Enter the FQDN for the vRealize Automation services"),
    [Parameter ()] 
        [ValidateNotNullOrEmpty()]
        [Alias('User')]
        [String]$Username = (Read-Host "Enter the username to authenticate with vRealize Automation"),
    [Parameter ()] 
        [ValidateNotNullOrEmpty()]
        [Alias('Pass')]
        [string]$Password = (Read-Host -MaskInput -Prompt "Enter the password to authenticate with vRealize Automation"),
    [Parameter ()] [String]$domain = (Read-Host -Prompt "Enter the domain or press enter to skip"),
    [Parameter ()] [switch]$SkipCertValidation=$false
)


Function Set-BasicAuthHeader {
    $base64AuthInfo = [Convert]::ToBase64String([Text.Encoding]::ASCII.GetBytes(("{0}:{1}" -f $username, $password))) # Create Basic Authentication Encoded Credentials
    $headers = @{"Accept" = "application/json" }
    $headers.Add("Authorization", "Basic $base64AuthInfo")
    $headers.Add("Content-Type", "application/json")
    $headers
}
Function Request-vRARefreshToken {
    Param (
        [String]$fqdn,
        [String]$username,
        [String]$password,
        [String]$domain,
        [switch]$SkipCertValidation
    )

    if ($SkipCertValidation) {
        add-type @"
        using System.Net;
        using System.Security.Cryptography.X509Certificates;
        public class TrustAllCertsPolicy : ICertificatePolicy {
            public bool CheckValidationResult(
            ServicePoint srvPoint, X509Certificate certificate,
            WebRequest request, int certificateProblem) {
            return true;
            }
        }
"@
    $AllProtocols = [System.Net.SecurityProtocolType]'Ssl3,Tls,Tls11,Tls12'
    [System.Net.ServicePointManager]::SecurityProtocol = $AllProtocols
    [System.Net.ServicePointManager]::CertificatePolicy = New-Object TrustAllCertsPolicy
    }

    $vraBasicHeaders = Set-BasicAuthHeader $username $password
    $Global:vraFqdn = $fqdn

    Try {
        $uri = "https://$vraFqdn/csp/gateway/am/api/login?access_token"
        if ($PsBoundParameters.ContainsKey("domain")) {
            $body = "{ ""username"":""$username"",""password"":""$password"",""domain"":""$tenant""}"
        }
        else {
            $body = "{ ""username"":""$username"",""password"":""$password""}"
        }

        if ($PSEdition -eq 'Core') {
            $vraResponse = Invoke-WebRequest -Method POST -Uri $uri -Headers $vraBasicHeaders -Body $body -SkipCertificateCheck # PS Core has -SkipCertificateCheck implemented, PowerShell 5.x does not.
        }
        else {
            $vraResponse = Invoke-WebRequest -Method POST -Uri $uri -Headers $vraBasicHeaders -Body $body
        }

        if ($vraResponse.StatusCode -eq 200) {
            $Global:vraHeaders = New-Object "System.Collections.Generic.Dictionary[[String],[String]]"
            $vraHeaders.Add("Accept", "application/json")
            $vraHeaders.Add("Content-Type", "application/json")
            $vraHeaders.Add("Authorization", "Bearer " + $vraResponse.Headers.'Csp-Auth-Token')
            Write-Output "Successfully connected to endpoint for vRealize Automation services: $vraFqdn"
            Write-Output "Generating token..."
            Write-Output "`n---------Refresh Token---------"
            ((Select-String -InputObject $vraResponse -Pattern '"refresh_token":') -Split ('"'))[3]
            Write-Output "-------------------------------`n"
            Write-Output "Saving environmental variables...`n"        
            $ENV:VRA_URL="https://$vRAFqdn"
            $ENV:VRA_REFRESH_TOKEN=((Select-String -InputObject $vraResponse -Pattern '"refresh_token":') -Split ('"'))[3]
            Write-Output "VRA_URL = $ENV:VRA_URL"
            Write-Output "VRA_REFRESH_TOKEN = $ENV:VRA_REFRESH_TOKEN"
        }
    }
    Catch {
        Write-Error $_.Exception.Message
    }
}

# Execute Functions
Request-vRARefreshToken -FQDN $FQDN -Username $Username -Password $Password -Domain $Domain -SkipCertValidation $SkipCertValidation
