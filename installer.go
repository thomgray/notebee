package main

import (
	"os"
	"path/filepath"

	"github.com/thomgray/notebee/config"
)

func install() {
	installDir := filepath.Join(config.GetAppConfig().HomeDir, ".notebee")

	if _, err := os.Stat(installDir); os.IsNotExist(err) {
		os.Mkdir(installDir, os.ModePerm)
	}
}
