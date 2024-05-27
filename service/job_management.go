package service

import "github.com/IceWhaleTech/CasaOS-JobManagement/codegen"

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

	m.jobMap[*job.ID] = job
}
