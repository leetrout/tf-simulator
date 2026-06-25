# terraform-provider-nimbus

A Terraform provider (built with
[terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework))
that manages resources in a running **tfsim** cloud via its HTTP API at
`/api/cloud/*`.

This is a standalone Go module (`github.com/leetrout/terraform-sim/cmd/terraform-provider-nimbus`)
kept separate from the root module so the heavy framework dependencies stay out
of the server binary.

## Resources

| Terraform type             | REST collection      | id prefix | writable attributes                  |
|----------------------------|----------------------|-----------|--------------------------------------|
| `nimbus_network`           | `networks`           | `net-`    | name, cidr_block, region             |
| `nimbus_subnet`            | `subnets`            | `subnet-` | name, network_id, cidr_block         |
| `nimbus_network_interface` | `network-interfaces` | `nic-`    | name, subnet_id, private_ip          |
| `nimbus_vm_instance`       | `vm-instances`       | `vm-`     | name, machine_type, nic_id           |
| `nimbus_storage_bucket`    | `storage-buckets`    | `bkt-`    | name, region                         |

Every resource also exposes computed `id` (string) and `serial` (number).
All resources implement Create / Read / Update (PATCH) / Delete and
`terraform import` (passthrough by `id`). Parent-reference attributes
(`network_id`, `subnet_id`, `nic_id`) force replacement when changed.

## Provider configuration

```hcl
provider "nimbus" {
  endpoint = "http://localhost:9321" # optional, this is the default
}
```

## Building

This repo uses a project-local Go toolchain (no system Go). From this
directory:

```bash
make build      # -> ./bin/terraform-provider-nimbus
```

The `Makefile` recipes source `../../.toolchain/env.sh` automatically. To run
the underlying commands by hand:

```bash
source ../../.toolchain/env.sh
go build -o bin/terraform-provider-nimbus .
go vet ./...
```

## Local install via `dev_overrides`

Terraform `dev_overrides` lets Terraform use a locally built binary directly,
bypassing the registry and `terraform init`. After `make build`, add the
following to `~/.terraformrc` (replace the path with your absolute checkout
path; `make install-dir` prints this for you):

```hcl
provider_installation {
  dev_overrides {
    "leetrout/nimbus" = "/Users/user/code/github.com/leetrout/terraform-sim/cmd/terraform-provider-nimbus/bin"
  }
  direct {}
}
```

The override value is the **directory** containing the
`terraform-provider-nimbus` binary, not the binary itself.

## Running the example

1. Start a `tfsim` server (listening on `http://localhost:9321`).
2. Build the provider and configure `dev_overrides` (above).
3. Apply the example topology:

```bash
cd examples
terraform plan     # no `init` needed under dev_overrides
terraform apply
```

The example (`examples/main.tf`) builds: one network -> two subnets
(public/private) -> a network interface -> a VM instance, plus two storage
buckets (`logs`, `assets`).

## Smoke test

`script/smoke.sh` runs a full `plan / apply / plan (idempotency) / destroy`
cycle against a running tfsim. It requires `terraform` (or set `TF_BIN=tofu`)
on `PATH`:

```bash
./script/smoke.sh
```
