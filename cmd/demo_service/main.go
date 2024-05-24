//go:generate bash -c "mkdir -p codegen/job_management && go run github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@v2.1.0 -generate types,client -package job_management ../../api/job_management/openapi.yaml > codegen/job_management/api.go"

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/vbauerster/mpb/v8"

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
		var message string
		if err := json.Unmarshal(response.Body, &baseResponse); err != nil {
			if len(response.Body) > 0 {
				message = fmt.Sprintf(" - %s", string(response.Body))
			}
		} else {
			if baseResponse.Message != nil {
				message = fmt.Sprintf(" - %s", *baseResponse.Message)
			}
		}

		fmt.Printf("%s%s\n", response.Status(), message)
		// os.Exit(1)
	}

	bars := mpb.NewWithContext(ctx)

	jobList := make([]*Job, 0)

	totalUnits := []int64{1000, 1000, 1000}
	unitTime := []time.Duration{100 * time.Millisecond, 100 * time.Millisecond, 100 * time.Millisecond}

	for i := 0; i < len(totalUnits); i++ {
		job := NewJob(totalUnits[i], unitTime[i])
		bar := bars.AddBar(job.Total())
		job.OnUnitCompletion(bar.Increment)
		jobList = append(jobList, job)
	}

	for _, job := range jobList {
		job.StartAsync(ctx)
	}

	bars.Wait()
}
