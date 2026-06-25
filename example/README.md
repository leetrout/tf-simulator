# statesim example

A runnable Terraform configuration for the Nimbus provider — the network → subnet
→ interface → vm topology plus two storage buckets.

## Run the teach-by-doing loop

1. Build the provider and add the `dev_overrides` snippet to `~/.terraformrc`
   (see [`../cmd/terraform-provider-nimbus/README.md`](../cmd/terraform-provider-nimbus/README.md)).

2. Start the server pointed at this directory, **without** the demo seed so it does
   not write a state file here:

   ```bash
   ./build/tfsim --work-dir ./example --seed empty
   ```

   (When the work dir contains `.tf` files, the demo seed is skipped automatically,
   so plain `./build/tfsim --work-dir ./example` works too.)

3. With `dev_overrides` configured you can skip `terraform init`. Apply:

   ```bash
   cd example
   terraform apply        # creates the resources in the fake cloud
   ```

   The UI graph fills in and every node shows `in_sync`.

4. Load a scenario in the UI (or `curl -X POST .../api/sim/scenarios/phantom/load`)
   to manufacture drift/phantom/etc., then run the suggested fix here and watch the
   dashboard reconcile.

> If you previously ran with the default `--seed demo` and got a stray
> `terraform.tfstate` in this directory with no plan, just delete it:
> `rm terraform.tfstate` — it was the no-Terraform demo state.
