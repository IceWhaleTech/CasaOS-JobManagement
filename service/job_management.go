package service

import (
	"sync"

	"github.com/IceWhaleTech/CasaOS-JobManagement/codegen"
	"github.com/samber/lo"
)

type JobManagement struct {
	jobMap    *sync.Map // TODO: use a persistent storage like SQLite
	nextJobID codegen.JobID
}

func NewJobManagement() *JobManagement {
	return &JobManagement{
		jobMap: &sync.Map{},
	}
}

func (m *JobManagement) GetJobMap() map[codegen.JobID]*codegen.Job {
	jobMap := map[codegen.JobID]*codegen.Job{}

	m.jobMap.Range(func(key, value any) bool {
		jobMap[key.(codegen.JobID)] = value.(*codegen.Job)
		return true
	})

	return jobMap
}

func (m *JobManagement) GetJob(jobID codegen.JobID) *codegen.Job {
	if job, ok := m.jobMap.Load(jobID); ok {
		return job.(*codegen.Job)
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

	m.jobMap.Store(*job.ID, job)

	m.nextJobID++
}

func (m *JobManagement) UpdateJobStatus(jobID codegen.JobID, jobStatus *codegen.JobStatus) {
	if job, ok := m.jobMap.Load(jobID); ok {
		job.(*codegen.Job).Status = jobStatus
		m.jobMap.Store(jobID, job)
	}
}

func (m *JobManagement) UpdateJobPriority(jobID codegen.JobID, jobPriority *codegen.JobPriority) {
	if job, ok := m.jobMap.Load(jobID); ok {
		job.(*codegen.Job).Priority = jobPriority
		m.jobMap.Store(jobID, job)
	}
}
