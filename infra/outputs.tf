output "dot_dev_ip" {
  value = linode_instance.dot_dev.ip_address
}

output "dot_dev_private_key" {
  value     = tls_private_key.dot_dev.private_key_openssh
  sensitive = true
}
