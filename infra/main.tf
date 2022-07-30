terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "4.26.0"
    }
  }
}

provider "google" {
  credentials = var.google_credentials
  project     = var.google_project
  region      = var.google_region
  zone        = var.google_zone
}

resource "google_service_account" "dot_dev" {
  account_id = "dot-dev-${var.google_region}"
}

resource "google_project_iam_binding" "instance-admin" {
  project = var.google_project
  role    = "roles/compute.instanceAdmin"
  members = ["serviceAccount:${google_service_account.dot_dev.email}"]
}

resource "google_project_iam_binding" "service-account" {
  project = var.google_project
  role    = "roles/iam.serviceAccountUser"
  members = ["serviceAccount:${google_service_account.dot_dev.email}"]
}

resource "google_service_account_key" "dot_dev" {
  service_account_id = google_service_account.dot_dev.name
}

resource "google_compute_address" "dot_dev" {
  name         = "dot-dev-ip"
  address_type = "EXTERNAL"
}

resource "google_compute_instance" "dot_dev" {
  name                      = "dot-dev-server"
  machine_type              = "e2-micro"
  zone                      = var.google_zone
  tags                      = ["web"]
  allow_stopping_for_update = true

  boot_disk {
    initialize_params {
      image = "cos-cloud/cos-stable"
    }
  }

  network_interface {
    network = "default"
    access_config {
      nat_ip = google_compute_address.dot_dev.address
    }
  }

  scheduling {
    automatic_restart   = true
    preemptible         = false
    on_host_maintenance = "MIGRATE"
  }

  service_account {
    email  = google_service_account.dot_dev.email
    scopes = ["cloud-platform", "compute-rw"]
  }

  metadata_startup_script = <<-EOT
  #!/bin/env bash

  cat > /etc/systemd/system/dot_dev.service <<-EOF
  [Unit]
  Description=aspenjames.dev web service
  After=docker.service
  Requires=docker.service

  [Service]
  TimeoutStartSec=0
  Restart=always
  ExecStartPre=-/usr/bin/docker pull aspenjames/aspenjames.dev:latest
  ExecStart=-/usr/bin/docker run --rm --name dot_dev -v /home/runner:/usr/src/content -p 80:3030 -t aspenjames/aspenjames.dev:latest
  ExecStop=/usr/bin/docker stop dot_dev
  EOF

  systemctl daemon-reload
  systemctl enable dot_dev
  systemctl start dot_dev
  EOT
}

resource "google_compute_firewall" "dot_dev" {
  name          = "dot-dev-web-traffic"
  network       = "default"
  source_ranges = ["0.0.0.0/0"]
  target_tags   = ["web"]

  allow {
    protocol = "icmp"
  }

  allow {
    protocol = "tcp"
    ports    = ["80"]
  }

  allow {
    protocol = "udp"
    ports    = ["80"]
  }
}

