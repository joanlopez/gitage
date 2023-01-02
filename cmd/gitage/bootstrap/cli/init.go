package cli

import (
	"github.com/spf13/cobra"

	"github.com/joanlopez/gitage"
)

func (c *CLI) initCmd() *cobra.Command {
	if c.init == nil {
		c.init = &cobra.Command{
			Use:   "init",
			Short: "Initialize a new Gitage repository",
			Long: `init is for initializing a new Gitage repository.
It creates the .gitage directory with the subsequent files.`,
			Args: cobra.ExactArgs(0),
			RunE: func(cmd *cobra.Command, args []string) error {
				return gitage.Init(c.ctx, c.fs, c.path, c.recipients...)
			},
		}
	}

	return c.init
}
