terraform {
  required_version = ">= 1.5.0"
  required_providers {
    docker = {
      source  = "kreuzwerker/docker"
      version = "~> 3.0"
    }
  }
}

variable "image" {
  type        = string
  description = "Pre-built ANGEL control-plane image."
  default     = "angel-control-plane:latest"
}

resource "docker_network" "angel_lab" {
  name = "angel-lab"
}

resource "docker_container" "control_plane" {
  name  = "angel-control-plane"
  image = var.image

  networks_advanced {
    name = docker_network.angel_lab.name
  }

  ports {
    internal = 8000
    external = 8000
    ip       = "127.0.0.1"
  }

  read_only     = true
  cap_drop      = ["ALL"]
  security_opts = ["no-new-privileges:true"]
  tmpfs = {
    "/tmp" = "rw,noexec,nosuid,size=16m"
  }
  restart = "unless-stopped"
}
