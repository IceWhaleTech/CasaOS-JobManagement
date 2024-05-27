package service

import (
	"github.com/IceWhaleTech/CasaOS-JobManagement/codegen"
	"github.com/samber/lo"
)

type JobManagement struct {
	jobMap    map[codegen.JobID]*codegen.Job // TODO: use a persistent storage like SQLite
	nextJobID codegen.JobID
}

func NewJobManagement() *JobManagement {
	return &JobManagement{
		jobMap: make(map[codegen.JobID]*codegen.Job, 0),
	}
}

func (m *JobManagement) GetJobMap() map[codegen.JobID]*codegen.Job {
	return m.jobMap
}

func (m *JobManagement) CreateJob(job *codegen.Job) {
	if job == nil {
		return
	}

	nextJobID := m.nextJobID // copy by value
	job.ID = &nextJobID      // copy by reference

	if job.Status == nil {
		status := codegen.JobStatus{}
		job.Status = &status
	}

	if job.Status.Status == "" {
		job.Status.Status = codegen.Running
	}

	if job.Status.Progress == nil {
		job.Status.Progress = lo.ToPtr(0)
	}

	m.jobMap[*job.ID] = job

	m.nextJobID++
}
