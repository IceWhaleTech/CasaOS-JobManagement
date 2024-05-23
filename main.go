//go:generate bash -c "mkdir -p codegen && go run github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@v2.1.0 -generate types,server,spec -package codegen api/job_management/openapi.yaml > codegen/job_management_api.go"

package main

import (
	_ "embed"
	"flag"
	"fmt"
	"os"

	"github.com/IceWhaleTech/CasaOS-Common/utils/logger"
	"github.com/IceWhaleTech/CasaOS-JobManager/common"
	"github.com/IceWhaleTech/CasaOS-JobManager/config"
	"github.com/IceWhaleTech/CasaOS-JobManager/service"
)

var (
	commit = "private build"
	date   = "private build"

	//go:embed api/index.html
	_docHTML string

	//go:embed api/job_management/openapi.yaml
	_docYAML string

	//go:embed build/sysroot/etc/casaos/job-management.conf.sample
	_confSample string
)

func main() {
	// parse arguments and intialize
	{
		configFlag := flag.String("c", "", "config file path")
		versionFlag := flag.Bool("v", false, "version")

		flag.Parse()

		if *versionFlag {
			fmt.Printf("version: %s\n", common.Version)
			fmt.Printf("git commit: %s\n", commit)
			fmt.Printf("build date: %s\n", date)

			os.Exit(0)
		}

		fmt.Printf("git commit: %s\n", commit)
		fmt.Printf("build date: %s\n", date)

		config.InitSetup(*configFlag, _confSample)

		logger.LogInit(config.AppInfo.LogPath, config.AppInfo.LogSaveName, config.AppInfo.LogFileExt)

		service.Initialize(config.CommonInfo.RuntimePath)
	}
}
