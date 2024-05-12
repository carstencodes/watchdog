package app

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/carstencodes/watchdog/internal/lib/common"
	"github.com/carstencodes/watchdog/internal/lib/config"
)

func initOs(fileProvider config.FileProvider) error {
	var err error = nil
	err = fileProvider.AddSearchPath("/etc")
	if err != nil {
		return err
	}
	err = fileProvider.AddSearchPath("/etc/" + common.ApplicationInfo().Name())
	if err != nil {
		return err
	}
	err = fileProvider.AddSearchPath("/etc/opt/" + common.ApplicationInfo().Name())
	if err != nil {
		return err
	}
	err = fileProvider.AddSearchPath("/opt/etc/" + common.ApplicationInfo().Name())
	if err != nil {
		return err
	}
	err = fileProvider.AddSearchPath(os.Getenv("HOME") + "/." + common.ApplicationInfo().Name() + "/")
	if err != nil {
		return err
	}
	err = fileProvider.AddSearchPath(os.Getenv("HOME") + "/.config/" + common.ApplicationInfo().Name() + "/")
	if err != nil {
		return err
	}
	err = fileProvider.AddSearchPath(os.Getenv("HOME") + "/.local/etc/" + common.ApplicationInfo().Name() + "/")
	if err != nil {
		return err
	}

	xdgDataDirs := os.Getenv("XDG_DATA_DIRS")
	xdgDataDirValues := strings.Split(xdgDataDirs, ":")
	for _, xdgDataDir := range xdgDataDirValues {
		err = fileProvider.AddSearchPath(filepath.Join(xdgDataDir, common.ApplicationInfo().Name()))
		if err != nil {
			return err
		}
	}

	return nil
}
