package run

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"code-intelligence.com/cifuzz/internal/build/cmake"
	"code-intelligence.com/cifuzz/internal/config"
)

// TODO: The reload command allows to reload the fuzz test names used
//       for autocompletion from the cmake config. It's only meant as a
//       temporary solution until we find a better solution.
type reloadCmd struct {
	*cobra.Command

	projectDir string
	config     *config.Config
}

func New(config *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reload [flags]",
		Short: "Reload fuzz test metadata",
		// TODO: Write long description
		Long: "",
		Args: cobra.NoArgs,
		RunE: func(c *cobra.Command, args []string) error {
			cmd := reloadCmd{Command: c, config: config}
			return cmd.run()
		},
	}
	return cmd
}

func (c *reloadCmd) run() error {
	if c.config.BuildSystem == config.BuildSystemCMake {
		return cmake.Reload(c.projectDir, c.OutOrStdout(), c.ErrOrStderr())
	} else if c.config.BuildSystem == config.BuildSystemUnknown {
		// Nothing to reload for unknown build system
		return nil
	} else {
		return errors.Errorf("Unsupported build system \"%s\"", c.config.BuildSystem)
	}
}
