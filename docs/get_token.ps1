# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE
# WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
# COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR
# OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

<#
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

    .EXAMPLE
    .\get_token.ps1
#>

Function Set-BasicAuthHeader {
    $base64AuthInfo = [Convert]::ToBase64String([Text.Encoding]::ASCII.GetBytes(("{0}:{1}" -f $username, $password))) # Create Basic Authentication Encoded Credentials
    $headers = @{"Accept" = "application/json" }
    $headers.Add("Authorization", "Basic $base64AuthInfo")
    $headers.Add("Content-Type", "application/json")
    $headers
}
Function Request-vRARefreshToken {
    Param (
        [Parameter (HelpMessage = "FQDN for the vRealize Automation.")] [String]$fqdn,
        [Parameter (HelpMessage = "Username to authenticate with vRealize Automation.")] [String]$username,
        [Parameter (HelpMessage = "Password to authenticate with vRealize Automation.")] [String]$password,
        [Parameter (HelpMessage = "Domain for the user or press enter to skip.")] [String]$domain,
        [Parameter (HelpMessage = "Skip certificate validation.")] [switch]$SkipCertValidation
    )

    if ($PSBoundParameters.Keys.Contains("fqdn")) { 
        Write-Host "FQDN variable found: $fqdn Skipping..."
    } else {
        $fqdn = Read-Host -Prompt "Enter the FQDN for the vRealize Automation services"
    }

    if ($PSBoundParameters.Keys.Contains("username")) { 
        Write-Host "Username variable found: $username. Skipping..."
    } else {
        $username = Read-Host -Prompt "Enter the username to authenticate with vRealize Automation"
    }

    if ($PSBoundParameters.Keys.Contains("password")) { 
        Write-Host "Password variable found. Skipping..."
    } else {
        $password = Read-Host -Prompt "Enter the password to authenticate with vRealize Automation"
    }

    if ($PSBoundParameters.Keys.Contains("domain")) { 
        Write-Host "Domain variable found: $domain. Skipping..."
    } else {
    }

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
Request-vRARefreshToken