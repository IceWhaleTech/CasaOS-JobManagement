//go:generate bash -c "mkdir -p codegen/job_management && go run github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@v2.1.0 -generate types,client -package job_management ../../api/job_management/openapi.yaml > codegen/job_management/api.go"

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/samber/lo"
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"

	"github.com/IceWhaleTech/CasaOS-JobManagement/cmd/demo_service/codegen/job_management"
)

const (
	sourceID              = "demo_service"
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

	jobManagement, err := job_management.NewClientWithResponses(url)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bars := mpb.NewWithContext(ctx)

	taskList := make([]*Task, 0)

	totalUnits := []int64{1000, 1000, 1000}
	unitTime := []time.Duration{100 * time.Millisecond, 100 * time.Millisecond, 100 * time.Millisecond}

	for i := 0; i < len(totalUnits); i++ {
		task := NewTask(totalUnits[i], unitTime[i])
		bar := bars.AddBar(task.totalUnits, mpb.PrependDecorators(
			decor.Any(func(decor.Statistics) string { return fmt.Sprintf("Job-%d", task.jobID) }),
		))
		task.onUnitCompletion = append(task.onUnitCompletion, func() {
			// upon completion of a unit, do following...
			bar.Increment()

			jobStatus := job_management.JobStatus{
				Status:   job_management.Running,
				Progress: lo.ToPtr(int((int64(i) / totalUnits[i]) * 100)),
			}

			response, _err := jobManagement.UpdateJobStatusWithResponse(ctx, task.jobID, jobStatus)
			if _err != nil {
				fmt.Println(_err.Error())
			} else if response.StatusCode() != http.StatusOK {
				printResponseMessage(response.HTTPResponse)
			}
		})
		taskList = append(taskList, task)
	}

	for _, task := range taskList {
		task.StartAsync(ctx, func(t *Task) {
			job := job_management.Job{
				SourceId: sourceID,
			}

			response, _err := jobManagement.CreateJobWithResponse(ctx, job)
			if _err != nil {
				fmt.Println(_err.Error())
				return
			}

			if response.StatusCode() != http.StatusOK {
				printResponseMessage(response.HTTPResponse)
				return
			}

			if response.JSON200 != nil && response.JSON200.Data != nil && response.JSON200.Data.ID != nil {
				t.jobID = *response.JSON200.Data.ID
			}
		})
	}

	bars.Wait()

	response, err := jobManagement.GetJobListWithResponse(ctx)
	if err != nil {
		panic(err)
	}

	if response.StatusCode() != http.StatusOK {
		printResponseMessage(response.HTTPResponse)
		// os.Exit(1)
	}

	if response.JSON200 == nil || response.JSON200.Data == nil {
		return
	}

	for _, job := range *response.JSON200.Data {
		fmt.Printf("Job-%d: %s (%s)\n", *job.ID, job.Status.Status, *job.Priority)
	}
}

func printResponseMessage(response *http.Response) {
	var message string
	var baseResponse job_management.BaseResponse

	body, err := io.ReadAll(response.Body)
	if err != nil {
		body = []byte{}
	}

	if err := json.Unmarshal(body, &baseResponse); err != nil {
		if len(body) > 0 {
			message = fmt.Sprintf(" - %s", string(body))
		}
	} else {
		if baseResponse.Message != nil {
			message = fmt.Sprintf(" - %s", *baseResponse.Message)
		}
	}

	fmt.Printf("%d - %s%s\n", response.StatusCode, response.Status, message)
}
