package service

import "github.com/IceWhaleTech/CasaOS-JobManagement/codegen"

type JobManagement struct {
	jobMap map[codegen.JobID]*codegen.Job
}

func NewJobManagement() *JobManagement {
	return &JobManagement{
		jobMap: make(map[codegen.JobID]*codegen.Job, 0),
	}
}

func (m *JobManagement) GetJobMap() map[codegen.JobID]*codegen.Job {
	return m.jobMap
}
