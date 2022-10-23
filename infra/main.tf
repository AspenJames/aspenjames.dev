terraform {
  required_providers {
    linode = {
      source  = "linode/linode"
      version = "1.29.4"
    }
  }
}

provider "linode" {
  token = var.linode_token
}

resource "tls_private_key" "dot_dev" {
  algorithm = "RSA"
  rsa_bits  = 4096
}

resource "linode_sshkey" "dot_dev" {
  label   = "dot_dev"
  ssh_key = chomp(tls_private_key.dot_dev.public_key_openssh)
}

resource "linode_stackscript" "dot_dev" {
  label       = "dot_dev"
  description = "Sets up aspenjames.dev Docker service"
  script      = <<EOS
#!/bin/bash

mkdir -p /usr/src/content

curl -fsSL https://get.docker.com | sh -
apt-get -q update
apt-get -q -y upgrade

cat > /etc/systemd/system/dot_dev.service <<-EOF
[Unit]
Description=aspenjames.dev web service
After=docker.service
Requires=docker.service

[Service]
TimeoutStartSec=0
Restart=always
ExecStartPre=-/usr/bin/docker pull aspenjames/aspenjames.dev:latest
ExecStart=-/usr/bin/docker run --rm --name dot_dev -v /usr/src/content:/usr/src/content -p 80:3030 -t aspenjames/aspenjames.dev:latest
ExecStop=/usr/bin/docker stop dot_dev
EOF

systemctl daemon-reload
systemctl enable dot_dev
systemctl start dot_dev
EOS
  images      = ["linode/ubuntu22.04"]
}

resource "linode_instance" "dot_dev" {
  label            = "web"
  image            = "linode/ubuntu22.04"
  region           = "us-west"
  type             = "g6-nanode-1"
  authorized_keys  = [linode_sshkey.dot_dev.ssh_key]
  watchdog_enabled = true
  stackscript_id   = linode_stackscript.dot_dev.id
}
