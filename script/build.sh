#!/usr/bin/env bash

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

cd ${SCRIPT_DIR}/../frontend
npm run build
npm run copy-to-go

cd ${SCRIPT_DIR}/..
mkdir -p build
go build -o build/tfsim cmd/tfsim/main.go
