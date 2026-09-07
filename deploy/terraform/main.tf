terraform {
  required_version = ">= 1.0"
  required_providers {
    digitalocean = {
      source  = "digitalocean/digitalocean"
      version = "~> 2.0"
    }
  }
}

provider "digitalocean" {
  token = var.do_token
}

resource "digitalocean_ssh_key" "default" {
  name       = "angel-ssh-key"
  public_key = file(var.ssh_public_key_path)
}

resource "digitalocean_vpc" "angel_vpc" {
  name     = "angel-vpc"
  region   = var.region
  ip_range = "10.0.0.0/16"
}

resource "digitalocean_droplet" "redirector" {
  image    = "ubuntu-24-04-x64"
  name     = "angel-redirector"
  region   = var.region
  size     = "s-1vcpu-1gb"
  vpc_uuid = digitalocean_vpc.angel_vpc.id
  ssh_keys = [digitalocean_ssh_key.default.id]

  tags = ["redirector", "c2-frontend"]

  user_data = <<-EOF
    #!/bin/bash
    apt-get update
    apt-get install -y nginx wireguard
    systemctl enable nginx
    systemctl start nginx
  EOF
}

resource "digitalocean_droplet" "teamserver" {
  image    = "ubuntu-24-04-x64"
  name     = "angel-teamserver"
  region   = var.region
  size     = "s-2vcpu-2gb"
  vpc_uuid = digitalocean_vpc.angel_vpc.id
  ssh_keys = [digitalocean_ssh_key.default.id]

  tags = ["teamserver", "c2-backend"]

  user_data = <<-EOF
    #!/bin/bash
    apt-get update
    apt-get install -y wireguard
  EOF
}

resource "digitalocean_droplet" "opsec" {
  image    = "ubuntu-24-04-x64"
  name     = "angel-opsec"
  region   = var.region
  size     = "s-1vcpu-1gb"
  vpc_uuid = digitalocean_vpc.angel_vpc.id
  ssh_keys = [digitalocean_ssh_key.default.id]

  tags = ["opsec", "recon"]
}

resource "digitalocean_droplet" "logging" {
  image    = "ubuntu-24-04-x64"
  name     = "angel-logging"
  region   = var.region
  size     = "s-1vcpu-1gb"
  vpc_uuid = digitalocean_vpc.angel_vpc.id
  ssh_keys = [digitalocean_ssh_key.default.id]

  tags = ["logging", "monitoring"]
}

output "redirector_ip" {
  value = digitalocean_droplet.redirector.ipv4_address
}

output "teamserver_ip" {
  value = digitalocean_droplet.teamserver.ipv4_address
}

output "opsec_ip" {
  value = digitalocean_droplet.opsec.ipv4_address
}

output "logging_ip" {
  value = digitalocean_droplet.logging.ipv4_address
}
