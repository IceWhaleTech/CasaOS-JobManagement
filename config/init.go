package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/IceWhaleTech/CasaOS-Common/utils/constants"
	"github.com/IceWhaleTech/CasaOS-JobManagement/common"
	"github.com/IceWhaleTech/CasaOS-JobManagement/model"
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
	ConfigFilePath = ModManagementConfigFilePath
	if len(config) > 0 {
		ConfigFilePath = config
	}

	// create default config file if not exist
	if _, err := os.Stat(ConfigFilePath); os.IsNotExist(err) {
		fmt.Println("config file not exist, create it")
		// create config file
		file, err := os.Create(ConfigFilePath)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		// write default config
		_, err = file.WriteString(sample)
		if err != nil {
			panic(err)
		}
	}

	var err error

	Cfg, err = ini.LoadSources(ini.LoadOptions{Insensitive: true, AllowShadows: true}, ConfigFilePath)
	if err != nil {
		panic(err)
	}

	mapTo("common", CommonInfo)
	mapTo("app", AppInfo)
}

func SaveSetup() error {
	reflectFrom("common", CommonInfo)
	reflectFrom("app", AppInfo)

	return Cfg.SaveTo(ConfigFilePath)
}

func mapTo(section string, v any) {
	err := Cfg.Section(section).MapTo(v)
	if err != nil {
		log.Fatalf("Cfg.MapTo %s err: %v", section, err)
	}
}

func reflectFrom(section string, v any) {
	err := Cfg.Section(section).ReflectFrom(v)
	if err != nil {
		log.Fatalf("Cfg.ReflectFrom %s err: %v", section, err)
	}
}
