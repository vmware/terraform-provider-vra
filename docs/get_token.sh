#!/bin/bash

# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE
# WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
# COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR
# OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

# Generates and returns a `refresh_token` from VMware Aria Automation for use by the Terraform provider.
#
#        terraform {
#            required_providers {
#                vra = {
#                    source  = "vmware/vra"
#                    version = ">= x.y.z"
#            }
#        }
#            required_version = ">= 0.13"
#        }
#
#        provider "vra" {
#            url           = "https://cloud.example.com"
#            refresh_token = "mx7w9**********************zB3UC"
#            insecure      = false
#        }
#
# Sets environment variables for `VRA_REFRESH_TOKEN` and `VRA_URL` for use by the Terraform provider.

### Check for an installtion of jq. ###

if ! [ -x "$(command -v jq)" ]; then
	echo -e "\nThe jq utility is missing. See https://stedolan.github.io/jq/ for installation instructions.\n"
	exit 1
fi

### Check for an existing endpoint value. ###

if [[ -v VRA_URL || -v fqdn ]]; then
	echo -e "\nFQDN variable found: $fqdn. Skipping...\n"
	export VRA_URL="https://$fqdn"
else
	echo -e "\nEnter the FQDN for the VMware Aria Automation services:"
	read fqdn
	export VRA_URL="https://$fqdn"
fi

### Check for an existing username value. ###

if [[ -v username ]]; then
	echo -e "\nUsername variable found: $username. Skipping...\n"
else
	echo -e "\nEnter the username to authenticate with VMware Aria Automation:"
	read username
fi

### Check for an existing password value. ###

if [[ -v password ]]; then
	echo -e "\nPassword variable found. Skipping...\n"
else
	echo -e "\nEnter the password to authenticate with VMware Aria Automation:"
	read -s password
fi

### Check for an a existing domain value. ###

if [[ -v domain ]]; then
	echo -e "\nDomain variable found: $domain. Skipping...\n"
else
	echo -e "\nEnter the domain or press enter to skip:"
	read domain
fi

### Generate the refresh token. ###

echo -e "\nGenerating Refresh Token..."
if [[ $domain == "" ]]; then
	export VRA_REFRESH_TOKEN=$(curl -k -X POST \
		"$VRA_URL/csp/gateway/am/api/login?access_token" \
		-H 'Content-Type: application/json' \
		-s \
		-d '{
		"username": "'"$username"'",
		"password": "'"$password"'"
		}' | jq -r .refresh_token)
else
	export VRA_REFRESH_TOKEN=$(curl -k -X POST \
		"$VRA_URL/csp/gateway/am/api/login?access_token" \
		-H 'Content-Type: application/json' \
		-s \
		-d '{
		"username": "'"$username"'",
		"password": "'"$password"'",
		"domain": "'"$domain"'"
		}' | jq -r .refresh_token)
fi

echo ""
echo "----------Refresh Token----------"
echo $VRA_REFRESH_TOKEN
echo "---------------------------------"
echo ""
echo "Environmental variables..."
echo ""
echo "VRA_URL = " $VRA_URL
echo "VRA_REFRESH_TOKEN = " $VRA_REFRESH_TOKEN

### Clear the password value. ###
unset password
