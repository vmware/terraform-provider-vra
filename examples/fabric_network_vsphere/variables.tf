variable "refresh_token" {
  description = "API Refresh token"
}

variable "url" {
  description = "URL to access the vRA(C) instance"
}

variable "cidr" {
  description = "CIDR notation is a compact representation of an IP address and its associated routing prefix"
}

variable "gateway" {
  description = "default IPv4 Gateway"
}

variable "domain" {
  description = "Domain name for the machine. The domain name is passed to the vSphere machine customization spec."
}


