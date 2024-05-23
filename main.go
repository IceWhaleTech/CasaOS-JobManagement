//go:generate bash -c "mkdir -p codegen && go run github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@v2.1.0 -generate types,server,spec -package codegen api/job_management/openapi.yaml > codegen/job_management_api.go"

package main

import (
	_ "embed"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/IceWhaleTech/CasaOS-Common/model"
	"github.com/IceWhaleTech/CasaOS-Common/utils/file"
	util_http "github.com/IceWhaleTech/CasaOS-Common/utils/http"
	"github.com/IceWhaleTech/CasaOS-Common/utils/logger"
	"github.com/IceWhaleTech/CasaOS-JobManagement/common"
	"github.com/IceWhaleTech/CasaOS-JobManagement/config"
	"github.com/IceWhaleTech/CasaOS-JobManagement/route"
	"github.com/IceWhaleTech/CasaOS-JobManagement/service"
	"github.com/coreos/go-systemd/daemon"

	"go.uber.org/zap"
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

	// setup listener
	listener, _err := net.Listen("tcp", net.JoinHostPort("127.0.0.1", "0"))
	if _err != nil { // use _err to avoid shadowing the err variables below.
		panic(_err)
	}

	// write address to file
	{
		urlFilePath := filepath.Join(config.CommonInfo.RuntimePath, "job-management.url")
		if err := file.CreateFileAndWriteContent(urlFilePath, "http://"+listener.Addr().String()); err != nil {
			logger.Error("error when creating address file", zap.Error(err),
				zap.Any("address", listener.Addr().String()),
				zap.Any("filepath", urlFilePath),
			)
		}
	}

	// initialize routers and register at gateway
	{
		apiPaths := []string{
			route.APIPath,
			route.DocPath,
		}

		for _, apiPath := range apiPaths {
			if err := service.MyService.Gateway().CreateRoute(&model.Route{
				Path:   apiPath,
				Target: "http://" + listener.Addr().String(),
			}); err != nil {
				panic(err)
			}
		}
	}

	router := route.GetRouter()
	docRouter := route.GetDocRouter(_docHTML, _docYAML)

	mux := &util_http.HandlerMultiplexer{
		HandlerMap: map[string]http.Handler{
			"v2":  router,
			"doc": docRouter,
		},
	}

	// notify systemd that we are ready
	{
		if supported, err := daemon.SdNotify(false, daemon.SdNotifyReady); err != nil {
			logger.Error("Failed to notify systemd that the service is ready", zap.Any("error", err))
		} else if supported {
			logger.Info("Notified systemd that the service is ready")
		} else {
			logger.Info("This process is not running as a systemd service.")
		}

		logger.Info("Module management service is listening...", zap.Any("address", listener.Addr().String()))
	}

	s := &http.Server{
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second, // fix G112: Potential slowloris attack (see https://github.com/securego/gosec)
	}

	_ = s.Serve(listener) // not using http.serve() to fix G114: Use of net/http serve function that has no support for setting timeouts (see https://github.com/securego/gosec)
	logger.Info("Job management service stopped")
}
