package cli

import (
	"github.com/spf13/cobra"

	"github.com/joanlopez/gitage"
)

func (c *CLI) initCmd() *cobra.Command {
	if c.init == nil {
		c.init = c.command(
			"init",
			"Initialize a new Gitage repository",
			`init is for initializing a new Gitage repository.
It creates the .gitage directory with the subsequent files.`,
		)

		// Set args
		c.init.Args = cobra.ExactArgs(0)

		// Set flags
		c.init.Flags().StringArrayVarP(&c.recipients, "recipient", "r", nil, "recipients to encrypt the repository")

		// Set run fn
		c.init.RunE = func(cmd *cobra.Command, args []string) error {
			return gitage.Init(c.ctx, c.fs, c.path, c.recipients...)
		}
	}

	return c.init
}
