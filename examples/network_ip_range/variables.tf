variable "refresh_token" {
    type = string
    description = "vRA Refresh Token"
}

variable "url" {
    type = string
    description = "vra url"
}

variable "cloud_account" {
    type = string
    description = "name of the cloud account"
}

variable "subnet_name" {
    type = string
    description = "name of the subnet to add ip range to"
}

variable "start_ip" {
    type = string
    description =  "starting ip for the range"
}

variable "end_ip" {
    type = string    
    description = "last ip of the range"
}

variable "ip_version" {
    type = string
    description = "ip version"
    default = "IPv4"
}