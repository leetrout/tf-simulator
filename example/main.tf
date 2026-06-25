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
# Resource `name` attributes are chosen so the built-in scenarios line up after a
# real `terraform apply` (the "drift" scenario targets the subnet named "private",
# "moved" targets the interface named "web0", "phantom" targets the bucket
# named "assets").

resource "nimbus_network" "main" {
  name       = "main"
  cidr_block = "10.0.0.0/16"
  region     = "us-west-1"
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

resource "nimbus_network_interface" "web0" {
  name       = "web0"
  subnet_id  = nimbus_subnet.public.id
  private_ip = "10.0.1.20"
}

resource "nimbus_vm_instance" "web" {
  name         = "web"
  machine_type = "n2-standard-2"
  nic_id       = nimbus_network_interface.web0.id
}

# --- Storage ------------------------------------------------------------------

resource "nimbus_storage_bucket" "logs" {
  name   = "logs"
  region = "us-west-1"
}

resource "nimbus_storage_bucket" "assets" {
  name   = "assets"
  region = "us-west-1"
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
