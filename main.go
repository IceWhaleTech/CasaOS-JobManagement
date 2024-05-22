//go:generate bash -c "mkdir -p codegen && go run github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@v2.1.0 -generate types,server,spec -package codegen api/job_management/openapi.yaml > codegen/job_management_api.go"

package main

import (
	_ "embed"
)

func main() {
}
