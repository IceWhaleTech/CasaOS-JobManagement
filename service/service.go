package service

import (
	"github.com/IceWhaleTech/CasaOS-Common/external"
)

var MyService *Services

type Services struct {
	gateway external.ManagementService

	runtimePath string
}

func Initialize(runtimePath string) {
	MyService = &Services{
		runtimePath: runtimePath,
	}
}
