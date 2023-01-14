package cli

import (
	"github.com/spf13/cobra"

	"github.com/joanlopez/gitage"
)

func (c *CLI) unregisterCmd() *cobra.Command {
	if c.unregister == nil {
		c.unregister = &cobra.Command{
			Use:   "unregister",
			Short: "Unregisters recipient(s) from the repository",
			Args:  cobra.ExactArgs(0),
			RunE: func(cmd *cobra.Command, args []string) error {
				return gitage.Unregister(c.ctx, c.fs, c.path, c.recipients...)
			},
		}
	}

	return c.unregister
}
