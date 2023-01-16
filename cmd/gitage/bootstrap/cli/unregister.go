package cli

import (
	"github.com/spf13/cobra"

	"github.com/joanlopez/gitage"
)

func (c *CLI) unregisterCmd() *cobra.Command {
	if c.unregister == nil {
		c.unregister = c.command(
			"unregister",
			"Unregisters recipient(s) from the repository",
			"",
		)

		// Set args
		c.unregister.Args = cobra.ExactArgs(0)

		// Set flags
		c.unregister.Flags().StringArrayVarP(&c.recipients, "recipient", "r", nil, "recipients to encrypt the repository")
		if err := c.unregister.MarkFlagRequired("recipient"); err != nil {
			panic(err)
		}

		// Set run fn
		c.unregister.RunE = func(cmd *cobra.Command, args []string) error {
			return gitage.Unregister(c.ctx, c.fs, c.path, c.recipients...)
		}
	}

	return c.unregister
}
