package config

import (
	"path/filepath"

	"github.com/IceWhaleTech/CasaOS-Common/utils/constants"
	"github.com/IceWhaleTech/CasaOS-JobManager/common"
	"github.com/IceWhaleTech/CasaOS-JobManager/model"
	"gopkg.in/ini.v1"
)

var (
	ModManagementConfigFilePath = filepath.Join(constants.DefaultConfigPath, "mod-management.conf")

	Cfg            *ini.File
	ConfigFilePath string

	CommonInfo = &model.CommonModel{
		RuntimePath: constants.DefaultRuntimePath,
	}

	AppInfo = &model.APPModel{
		LogPath:     constants.DefaultLogPath,
		LogSaveName: common.ServiceName,
		LogFileExt:  "log",
	}
)

func InitSetup(config string, sample string) {
}
