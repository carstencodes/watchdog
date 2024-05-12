package config

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/carstencodes/watchdog/internal/lib/common"
)

var configFileCliArgs = flag.NewFlagSet("config", flag.ContinueOnError)

var defaultConfigFileName = common.ApplicationInfo().Name() + ".json"

var configFilePath = defaultConfigFileName

func init() {
	cwd, err := os.Getwd()
	if err != nil {
		configFilePath = filepath.Join(cwd, defaultConfigFileName)
	}

	configFileCliArgs.StringVar(&configFilePath, "config", configFilePath, "Path to configuration file")
	flag.StringVar(&configFilePath, "config", configFilePath, "Path to configuration file")
}

type FileProvider interface {
	Provider
	AddSearchPath(path string) error
}
