package service

import (
	"github.com/IceWhaleTech/CasaOS-Common/external"
)

var MyService *Services

type Services struct {
	gateway external.ManagementService

	jobManagement *JobManagement

	runtimePath string
}

func Initialize(runtimePath string) {
	MyService = &Services{
		runtimePath: runtimePath,
	}
}

func (s *Services) Gateway() external.ManagementService {
	if s.gateway == nil {
		gateway, err := external.NewManagementService(s.runtimePath)
		if err != nil && len(s.runtimePath) > 0 {
			panic(err)
		}

		s.gateway = gateway
	}

	return s.gateway
}

func (s *Services) JobManagement() *JobManagement {
	if s.jobManagement == nil {
		s.jobManagement = NewJobManagement()
	}

	return s.jobManagement
}
