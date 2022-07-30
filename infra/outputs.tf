output "dot_dev_ip" {
  value = google_compute_instance.dot_dev.network_interface[0].access_config[0].nat_ip
}

output "dot_dev_private_key" {
  value     = base64decode(google_service_account_key.dot_dev.private_key)
  sensitive = true
}

output "dot_dev_public_key" {
  value     = base64decode(google_service_account_key.dot_dev.public_key)
  sensitive = true
}
