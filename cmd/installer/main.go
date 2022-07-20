package main

import (
	"os"

	"github.com/spf13/cobra"

	installer_bundle "code-intelligence.com/cifuzz/installer-bundle"
	"code-intelligence.com/cifuzz/pkg/log"
	"code-intelligence.com/cifuzz/tools/install"
)

func main() {
	opts := &install.Options{}
	fs := &installer_bundle.Bundle

	cmd := &cobra.Command{
		Use:   "installer",
		Short: "Install cifuzz",
		Run: func(cmd *cobra.Command, args []string) {
			install.ExtractBundle(opts.InstallDir, fs)
		},
	}

	cmd.Flags().StringVarP(&opts.InstallDir, "install-dir", "i", "~/cifuzz", "The directory to install cifuzz in")

	err := cmd.Execute()
	if err != nil {
		log.Error(err, err.Error())
		os.Exit(1)
	}
}
