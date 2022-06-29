variable "google_credentials" {
  type        = string
  description = "Google Cloud Project credentials JSON"
  sensitive   = true
}

variable "google_project" {
  type        = string
  default     = "nimble-granite-354015"
  description = "Google Cloud Project id"
}

variable "google_region" {
  type        = string
  default     = "us-west1"
  description = "Google Cloud region"
}

variable "google_zone" {
  type        = string
  default     = "us-west1-b"
  description = "Google Cloud zone"
}
