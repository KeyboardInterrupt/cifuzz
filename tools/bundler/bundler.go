package main

import (
	"flag"
	"os"

	"code-intelligence.com/cifuzz/pkg/log"
	"code-intelligence.com/cifuzz/tools/install"
)

func main() {
	var version string
	flag.StringVar(&version, "version", "v0.0.1", "the version used for cifuzz")

	opts := &install.Options{}
	installer, err := install.NewInstaller(opts)
	checkError(err)

	// TODO: compile for different os
	err = installer.InstallMinijail()
	checkError(err)
	err = installer.InstallProcessWrapper()
	checkError(err)
	_, err = installer.CopyCMakeIntegration(installer.ShareDir())
	checkError(err)
	err = installer.InstallCIFuzz(version)
	checkError(err)
}

func checkError(err error) {
	if err != nil {
		log.Error(err, err.Error())
		os.Exit(1)
	}
}
