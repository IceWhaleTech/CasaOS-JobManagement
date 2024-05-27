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

func (m *JobManagement) GetJob(jobID codegen.JobID) *codegen.Job {
	if job, ok := m.jobMap[jobID]; ok {
		return job
	}

	return nil
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

func (m *JobManagement) UpdateJobStatus(jobID codegen.JobID, jobStatus *codegen.JobStatus) {
	if job, ok := m.jobMap[jobID]; ok {
		job.Status = jobStatus
	}
}

func (m *JobManagement) UpdateJobPriority(jobID codegen.JobID, jobPriority *codegen.JobPriority) {
	if job, ok := m.jobMap[jobID]; ok {
		job.Priority = jobPriority
	}
}
