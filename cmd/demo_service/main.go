//go:generate bash -c "mkdir -p codegen/job_management && go run github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@v2.1.0 -generate types,client -package job_management ../../api/job_management/openapi.yaml > codegen/job_management/api.go"

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path"

	"github.com/IceWhaleTech/CasaOS-JobManagement/cmd/demo_service/codegen/job_management"
)

const (
	basePathJobManagement = "v2/job_management"
)

func main() {
	if len(os.Args) < 2 {
		os.Args = append(os.Args, "http://localhost")
	}

	inputFlag := os.Args[1]

	if inputFlag == "-h" || inputFlag == "--help" {
		fmt.Printf("usage: %s http://<base URL>\t(default: http://localhost)\n", path.Base(os.Args[0]))

		os.Exit(0)
	}

	url := fmt.Sprintf("%s/%s", inputFlag, basePathJobManagement)

	jobManagementClient, err := job_management.NewClientWithResponses(url)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	response, err := jobManagementClient.GetJobListWithResponse(ctx)
	if err != nil {
		panic(err)
	}

	if response.StatusCode() != http.StatusOK {
		var baseResponse job_management.BaseResponse
		if err := json.Unmarshal(response.Body, &baseResponse); err != nil {
			fmt.Printf("%s - %s\n", response.Status(), response.Body)
			os.Exit(1)
		}

		var message string
		if baseResponse.Message != nil {
			message = fmt.Sprintf(" - %s", *baseResponse.Message)
		}

		fmt.Printf("%s%s\n", response.Status(), message)
		os.Exit(1)
	}
}
