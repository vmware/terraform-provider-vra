variable "url" {
  description = "URL for reaching the vRA Cloud/vRA8 instance"
}
variable "refresh_token" {
  description = "Refresh token"
}

variable "insecure" {
  description = "Whether to allow for self-signed certs set true for vRA8 on prem, false for SaaS"
}

variable "integration_id" {
  description = "ID of the Gitlab/Github integration that should be used to access the repository"
}

variable "content_source_name" {
  description = "Name of the new content source"
}

variable "project_name" {
  description = "Name of the new project"
}
