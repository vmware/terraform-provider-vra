#!/bin/bash
#
# Script to generate a refresh token for vRA8 on prem or vRA Cloud. 
# This will prompt for the following values if they are not already set:
#    username
# Sets environment variables for VRA_REFRESH_TOKEN and VRA_URL which can be consumed by the 
# TF provider more securely than leaving the token in cleartext. 
#
#
if  ! [ -x "$(command -v jq)" ]
then
	echo "The jq utility is missing. See https://stedolan.github.io/jq/ for instructions to get it"
	exit 1
fi

#Check for an already existing username value
if [[ ! -z "$username" ]] 
then
	echo "username variable found: $username"
else 
	echo "Please enter username to connect to vra with"
	read username
fi

#Check for an already existing password value
if [[ ! -z "$password" ]]
then
	echo "password variable found"
else
	echo "Please enter password to connect to vra with"
	read -s password
fi

#Check for an already existing LDAP/AD domain value
if [[ ! -z "$domain" ]]
then
	echo "Existing domain variable found: $domain"
else
	echo "Please enter domain to connect to vra with (for AD/LDAP users) or press Enter if you not want to use domain"
	read domain
fi

if [[ -z "$VRA_URL" ]]
then 
	echo "Please enter the hostname/fqdn of the VRA8 server/ or cloud identity server"
	read host
	export VRA_URL="https://$host"
fi

echo "Using $VRA_URL"
#use different json bodies with curl depending on whether or not a domain 
# was specified
echo "Getting Token"
if [[ $domain == "" ]]
then
  curlArgs=(
	-k -X POST
	"$VRA_URL/csp/gateway/am/api/login?access_token"
	-H 'Content-Type: application/json'
	-s
	-d '{ "username": "'"$username"'", "password": "'"$password"'" }'
  )
else
  curlArgs=(
	-k -X POST
	"$VRA_URL/csp/gateway/am/api/login?access_token"
	-H 'Content-Type: application/json'
	-s
	-d '{ "username": "'"$username"'", "password": "'"$password"'", "domain": "'"$domain"'" }'
  )
fi

export TF_VAR_VRA_REFRESH_TOKEN=$(curl "${curlArgs[@]}" | jq -r .refresh_token)

#clean up password 
unset password

echo "Refresh Token"
echo "----------------------------"
echo $TF_VAR_VRA_REFRESH_TOKEN
echo "----------------------------"
