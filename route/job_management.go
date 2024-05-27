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

	return ctx.JSON(http.StatusNotImplemented, struct{}{})
}

func (m *JobManagement) GetJob(ctx echo.Context, _ codegen.JobID) error {
	return ctx.JSON(http.StatusNotImplemented, struct{}{})
}

func (m *JobManagement) UpdateJobPriority(ctx echo.Context, _ codegen.JobID) error {
	return ctx.JSON(http.StatusNotImplemented, struct{}{})
}

func (m *JobManagement) GetJobStatus(ctx echo.Context, _ codegen.JobID) error {
	return ctx.JSON(http.StatusNotImplemented, struct{}{})
}

func (m *JobManagement) UpdateJobStatus(ctx echo.Context, _ codegen.JobID) error {
	return ctx.JSON(http.StatusNotImplemented, struct{}{})
}
