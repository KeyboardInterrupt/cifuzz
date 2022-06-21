package bundle

import (
	"github.com/spf13/cobra"

	"code-intelligence.com/cifuzz/internal/build/cmake"
	"code-intelligence.com/cifuzz/internal/config"
	"code-intelligence.com/cifuzz/pkg/artifact"
	"code-intelligence.com/cifuzz/pkg/dialog"
	"code-intelligence.com/cifuzz/pkg/log"
)

type bundleCmd struct {
	*cobra.Command

	config *config.Config
	opts   *bundleOpts
}

type bundleOpts struct {
	outputPath string
}

func New(config *config.Config) *cobra.Command {
	opts := &bundleOpts{}
	bundleCmd := &cobra.Command{
		Use:   "bundle",
		Short: "Generates a artifact bundle",
		Long:  "collects all data necessary for an artifact bundle and stores it to a an archive",
		Args:  cobra.NoArgs,
		RunE: func(c *cobra.Command, args []string) error {
			cmd := bundleCmd{Command: c, config: config, opts: opts}
			return cmd.run()
		},
	}

	// TODO add option for selecting single fuzz tests for the artifact
	bundleCmd.Flags().StringVarP(&opts.outputPath, "output", "o", "", "File path of artifact file")

	return bundleCmd
}

func (c *bundleCmd) run() (err error) {
	// TODO offer completion for directories
	if c.opts.outputPath == "" {
		c.opts.outputPath, err = dialog.Input(
			"Please enter filename",
			"artifact.tar.gz",
			c.InOrStdin(),
		)
		if err != nil {
			return err
		}
	}

	if c.config.BuildSystem == config.BuildSystemCMake {
		// TODO reload is not enought... we should build the fuzz test(s)
		if err := cmake.Reload(c.config.ProjectDir, c.OutOrStdout(), c.ErrOrStderr()); err != nil {
			return err
		}
	}

	metadata, err := artifact.NewMetadata(c.config)
	if err != nil {
		return err
	}

	if err := artifact.Create(c.opts.outputPath, metadata); err != nil {
		return err
	}

	log.Successf("Successfully created artifact: %s", c.opts.outputPath)
	return nil
}
