#!/bin/sh

export SALES_DB_DISABLE_TLS=true

go run ./cmd/sales-admin "$@"

