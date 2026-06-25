terraform {
  required_providers {
    nimbus = {
      source = "leetrout/nimbus"
    }
  }
}

provider "nimbus" {
  endpoint = "http://localhost:9321"
}

# --- Network topology ---------------------------------------------------------

resource "nimbus_network" "main" {
  name       = "main"
  cidr_block = "10.0.0.0/16"
  region     = "us-east-1"
}

resource "nimbus_subnet" "public" {
  name       = "public"
  network_id = nimbus_network.main.id
  cidr_block = "10.0.1.0/24"
}

resource "nimbus_subnet" "private" {
  name       = "private"
  network_id = nimbus_network.main.id
  cidr_block = "10.0.2.0/24"
}

resource "nimbus_network_interface" "web" {
  name       = "web-nic"
  subnet_id  = nimbus_subnet.public.id
  private_ip = "10.0.1.10"
}

resource "nimbus_vm_instance" "web" {
  name         = "web-1"
  machine_type = "small"
  nic_id       = nimbus_network_interface.web.id
}

# --- Storage ------------------------------------------------------------------

resource "nimbus_storage_bucket" "logs" {
  name   = "logs"
  region = "us-east-1"
}

resource "nimbus_storage_bucket" "assets" {
  name   = "assets"
  region = "us-east-1"
}

# --- Outputs ------------------------------------------------------------------

output "network_id" {
  value = nimbus_network.main.id
}

output "vm_id" {
  value = nimbus_vm_instance.web.id
}

output "bucket_ids" {
  value = [nimbus_storage_bucket.logs.id, nimbus_storage_bucket.assets.id]
}
