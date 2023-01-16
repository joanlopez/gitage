package cli

import (
	"github.com/spf13/cobra"

	"github.com/joanlopez/gitage"
)

func (c *CLI) registerCmd() *cobra.Command {
	if c.register == nil {
		c.register = c.command(
			"register",
			"Registers new recipient(s) to the repository",
			"",
		)

		// Set args
		c.register.Args = cobra.ExactArgs(0)

		// Set flags
		c.register.Flags().StringArrayVarP(&c.recipients, "recipient", "r", nil, "recipients to encrypt the repository")
		if err := c.register.MarkFlagRequired("recipient"); err != nil {
			panic(err)
		}

		// Set run fn
		c.register.RunE = func(cmd *cobra.Command, args []string) error {
			return gitage.Register(c.ctx, c.fs, c.path, c.recipients...)
		}
	}

	return c.register
}
