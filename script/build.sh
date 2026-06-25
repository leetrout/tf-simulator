#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
ROOT="${SCRIPT_DIR}/.."

# 1. Build the frontend and embed it into the server binary.
cd "${ROOT}/frontend"
npm run build
npm run copy-to-go

# 2. Build the server binary (serves cloud API + sim API + ws + embedded UI).
cd "${ROOT}"
mkdir -p build
go build -o build/tfsim ./cmd/tfsim

# 3. Build the standalone Terraform provider plugin (separate module).
cd "${ROOT}/cmd/terraform-provider-nimbus"
go build -o "${ROOT}/build/terraform-provider-nimbus" .

echo "built: build/tfsim, build/terraform-provider-nimbus"
