package app

import (
	"github.com/carstencodes/watchdog/internal/lib/config"
)

func initOs(fileProvider config.FileProvider) error {
	var err error = nil
	err = fileProvider.AddSearchPath(getEnvDir("ALLUSERSPROFILE"))
	if err != nil {
		return err
	}

	err = fileProvider.AddSearchPath(getEnvDir("USERPROFILE"))
	if err != nil {
		return err
	}

	err = fileProvider.AddSearchPath(getEnvDir("APPDATA"))
	if err != nil {
		return err
	}

	return nil
}
