package app

import (
	"os"
	"path/filepath"

	"github.com/carstencodes/watchdog/internal/lib/common"
	"github.com/carstencodes/watchdog/internal/lib/config"
)

func setupConfiguration() error {
	var err error = nil
	config.Sections, err = config.Providers.BuildDefinitions()
	if err != nil {
		return err
	}

	config.Values, err = config.Sections.BuildConfiguration()
	if err != nil {
		return err
	}

	if !config.Values.Parsed() {
		return config.Values.Parse()
	}

	return nil
}

func init() {
	err := initConfig()
	if err != nil {
		panic(err)
	}
}

func initConfig() error {
	var err error = nil
	var namespace = common.ApplicationInfo().Name()

	json := config.NewJsonFileConfigProvider(false)
	err = initOs(json)
	if err != nil {
		return err
	}
	env := config.NewEnvProvider(namespace)
	flags := config.NewFlagConfigProvider()

	err = config.Providers.AddProvider(json)
	if err != nil {
		return err
	}
	err = config.Providers.AddProvider(env)
	if err != nil {
		return err
	}
	err = config.Providers.AddProvider(flags)

	return err
}

func getEnvDir(variableName string) string {
	value := os.Getenv(variableName)
	if value == "" {
		return ""
	}

	return filepath.Join(value, common.ApplicationInfo().Name())
}
