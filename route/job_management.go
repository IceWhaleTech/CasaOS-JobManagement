package route

import (
	"net/http"

	"github.com/IceWhaleTech/CasaOS-JobManagement/codegen"
	"github.com/IceWhaleTech/CasaOS-JobManagement/service"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
)

func (m *JobManagement) GetJobList(ctx echo.Context) error {
	jobMap := service.MyService.JobManagement().GetJobMap()

	jobList := lo.MapToSlice(jobMap, func(_ codegen.JobID, v *codegen.Job) codegen.Job {
		return *v
	})

	return ctx.JSON(http.StatusOK, codegen.JobListOK{Data: &jobList})
}

func (m *JobManagement) CreateJob(ctx echo.Context) error {
	var job codegen.Job
	if err := ctx.Bind(&job); err != nil {
		message := err.Error()
		return ctx.JSON(http.StatusBadRequest, codegen.ResponseBadRequest{Message: &message})
	}

	service.MyService.JobManagement().CreateJob(&job)

	return ctx.JSON(http.StatusOK, codegen.JobOK{Data: &job})
}

func (m *JobManagement) GetJob(ctx echo.Context, jobID codegen.JobID) error {
	job := service.MyService.JobManagement().GetJob(jobID)

	if job == nil {
		message := "job not found"
		return ctx.JSON(http.StatusNotFound, codegen.ResponseNotFound{Message: &message})
	}

	return ctx.JSON(http.StatusOK, codegen.JobOK{Data: job})
}

func (m *JobManagement) UpdateJobPriority(ctx echo.Context, jobID codegen.JobID) error {
	job := service.MyService.JobManagement().GetJob(jobID)

	if job == nil {
		message := "job not found"
		return ctx.JSON(http.StatusNotFound, codegen.ResponseNotFound{Message: &message})
	}

	var jobPriority codegen.JobPriority
	if err := ctx.Bind(&jobPriority); err != nil {
		message := err.Error()
		return ctx.JSON(http.StatusBadRequest, codegen.ResponseBadRequest{Message: &message})
	}

	service.MyService.JobManagement().UpdateJobPriority(jobID, &jobPriority)

	job = service.MyService.JobManagement().GetJob(jobID)
	if job == nil {
		message := "job not found - how possible??"
		return ctx.JSON(http.StatusNotFound, codegen.ResponseNotFound{Message: &message})
	}

	return ctx.JSON(http.StatusOK, codegen.JobPriorityOK{Data: job.Priority})
}

func (m *JobManagement) GetJobStatus(ctx echo.Context, jobID codegen.JobID) error {
	job := service.MyService.JobManagement().GetJob(jobID)

	if job == nil {
		message := "job not found"
		return ctx.JSON(http.StatusNotFound, codegen.ResponseNotFound{Message: &message})
	}

	return ctx.JSON(http.StatusOK, codegen.JobStatusOK{Data: job.Status})
}

func (m *JobManagement) UpdateJobStatus(ctx echo.Context, jobID codegen.JobID) error {
	job := service.MyService.JobManagement().GetJob(jobID)

	if job == nil {
		message := "job not found"
		return ctx.JSON(http.StatusNotFound, codegen.ResponseNotFound{Message: &message})
	}

	var jobStatus codegen.JobStatus
	if err := ctx.Bind(&jobStatus); err != nil {
		message := err.Error()
		return ctx.JSON(http.StatusBadRequest, codegen.ResponseBadRequest{Message: &message})
	}

	service.MyService.JobManagement().UpdateJobStatus(jobID, &jobStatus)

	job = service.MyService.JobManagement().GetJob(jobID)
	if job == nil {
		message := "job not found - how possible??"
		return ctx.JSON(http.StatusNotFound, codegen.ResponseNotFound{Message: &message})
	}

	return ctx.JSON(http.StatusOK, codegen.JobStatusOK{Data: job.Status})
}
