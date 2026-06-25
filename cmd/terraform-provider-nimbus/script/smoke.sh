#!/usr/bin/env bash
#
# smoke.sh — manual end-to-end smoke test for the nimbus provider.
#
# Prerequisites:
#   1. A running tfsim server on http://localhost:9321 (serving /api/cloud/*).
#   2. The provider binary built (see ../Makefile: `make build`).
#   3. ~/.terraformrc configured with a dev_overrides block pointing at the
#      built binary directory (see ../README.md). With dev_overrides you do NOT
#      run `terraform init` — Terraform uses the override directly.
#   4. terraform (or tofu) on PATH.
#
# NOTE: terraform/tofu are not installed in the build sandbox, so this script is
# documented for manual execution against a real environment.
#
# Usage:
#   ./script/smoke.sh
#
set -euo pipefail

# terraform or tofu
TF="${TF_BIN:-terraform}"
EXAMPLE_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../examples" && pwd)"

echo "==> Using $TF in $EXAMPLE_DIR"
cd "$EXAMPLE_DIR"

# With dev_overrides, `init` will warn and is unnecessary; we still try it so
# this also works without dev_overrides (a published/local-mirror install).
echo "==> terraform init (ignored under dev_overrides)"
"$TF" init -input=false || true

echo "==> terraform plan"
"$TF" plan -input=false

echo "==> terraform apply"
"$TF" apply -input=false -auto-approve

echo "==> terraform plan (expect: no changes / idempotent)"
"$TF" plan -input=false -detailed-exitcode || {
  code=$?
  if [ "$code" -eq 2 ]; then
    echo "WARNING: plan reported changes after apply (not idempotent)"
  else
    exit "$code"
  fi
}

echo "==> outputs"
"$TF" output

echo "==> terraform destroy"
"$TF" destroy -input=false -auto-approve

echo "==> smoke test complete"
