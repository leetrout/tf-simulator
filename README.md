# statesim

A laptop-friendly tool for **learning Terraform by doing**. A single Go binary
starts a local web server, serves a dark two-pane UI, and exposes a **fake cloud
provider** ("Nimbus"). You write real HCL, run real `terraform` commands in your
own shell against Nimbus using a **real Terraform provider**, and the UI shows —
live — what Terraform thinks exists (`terraform.tfstate`) versus what actually
exists in the fake cloud, classifying every resource as **in_sync / drift /
phantom / untracked / moved**. Inject scenarios to manufacture each condition and
practice the commands that resolve it.

## Goals (from the original repo)

- A single binary to run a local server with a dummy API.
- A hands-on refresher on how to write a Terraform provider.

## Architecture

Everything ships in one server binary plus a standalone provider plugin:

```
internal/
  cloud/      Nimbus fake cloud: models, store, REST handlers under /api/cloud/*
  sim/        observer: tfstate parser, diff engine, graph, scenarios, /api/sim/*
  ws/         structured websocket event hub (/ws)
  webserver/  routing + embedded frontend (//go:embed)
cmd/
  tfsim/                       the server binary (cloud API + sim API + ws + UI)
  terraform-provider-nimbus/   the provider plugin (separate Go module)
frontend/                      Svelte 5 + SvelteKit + Tailwind + daisyUI
```

The simulator reads the cloud store **in-process**; only the provider talks to
the cloud over HTTP (Terraform runs it out-of-process). The simulator is an
observer and scenario injector — it never runs `terraform` itself.

### Resource model (Nimbus)

| Terraform type             | collection           | attributes                       | depends on |
|----------------------------|----------------------|----------------------------------|------------|
| `nimbus_network`           | `networks`           | name, cidr_block, region         | —          |
| `nimbus_subnet`            | `subnets`            | name, network_id, cidr_block     | network    |
| `nimbus_network_interface` | `network-interfaces` | name, subnet_id, private_ip      | subnet     |
| `nimbus_vm_instance`       | `vm-instances`       | name, machine_type, nic_id       | interface  |
| `nimbus_storage_bucket`    | `storage-buckets`    | name, region                     | —          |

Every resource has a server-assigned `id` (the join key with `terraform.tfstate`)
and a `serial` bumped on each mutation.

### Status taxonomy

- `in_sync` — present in both, attributes equal.
- `drift` — present in both, attributes differ.
- `phantom` — in state, absent from cloud (a plan would recreate it).
- `untracked` — in cloud, absent from state (needs `import`).
- `moved` — same cloud `id`, only the name differs (renamed; needs a `moved` block).

## Run it

```bash
# build the server (embeds the frontend) and start it
./script/build.sh
./build/tfsim --work-dir ./example          # opens the UI at http://localhost:9321
```

Flags: `--addr` (`:9321`), `--work-dir` (`.`), `--seed demo|empty`, `--no-open`.

With `--seed demo` (default) the cloud is seeded with a seven-resource baseline
and a matching `terraform.tfstate` is written into the work dir, so the full
diff/scenario loop is demoable **without** installing Terraform.

## The teach-by-doing loop (with a real provider)

1. Build the provider and point Terraform at it with a `dev_overrides` block —
   see [`cmd/terraform-provider-nimbus/README.md`](cmd/terraform-provider-nimbus/README.md).
2. Start `tfsim --work-dir ./example --seed empty`.
3. In `./example`, `terraform init && terraform apply` — five resources appear in
   the cloud; the UI graph fills in; all nodes are `in_sync`.
4. Load a scenario (UI or `POST /api/sim/scenarios/phantom/load`) to manufacture a
   condition, then run the suggested fix in your shell and watch the UI reconcile.

## API surface

Fake cloud (REST, one collection per type):

```
GET|POST          /api/cloud/{collection}
GET|PATCH|PUT|DELETE  /api/cloud/{collection}/{id}
```

Simulator (read + scenario actions):

```
GET  /api/sim/snapshot            { cloud[], state[], status[], counts, meta }
GET  /api/sim/graph?view=cloud|state
GET  /api/sim/status
GET  /api/sim/settings
GET  /api/sim/scenarios
POST /api/sim/scenarios/{id}/load
POST /api/sim/reset
GET  /ws                          structured event stream
```

See [`requests/nimbus.http`](requests/nimbus.http) for runnable examples.

## Development

```bash
# Go
go test ./...
go run ./cmd/tfsim --work-dir /tmp/demo

# frontend (from frontend/)
npm run check && npm run lint && npm run build
```

> Sandbox note: this repo includes `.toolchain/env.sh`, which is only used when a
> system Go toolchain is unavailable. Ignore it if you have Go installed normally.
