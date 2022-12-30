package cli

import (
	"github.com/joanlopez/gitage"
	"github.com/spf13/cobra"
)

func (c *CLI) registerCmd() *cobra.Command {
	if c.register == nil {
		c.register = &cobra.Command{
			Use:   "register",
			Short: "Registers new recipient(s) to the repository",
			Args:  cobra.ExactArgs(0),
			RunE: func(cmd *cobra.Command, args []string) error {
				return gitage.Register(c.ctx, c.fs, c.path, c.recipients...)
			},
		}
	}

	return c.register
}
