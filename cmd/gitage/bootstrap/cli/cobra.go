package cli

import (
	"github.com/spf13/cobra"
)

func (c *CLI) command(use, short, long string) *cobra.Command {
	if len(long) == 0 {
		long = short
	}

	cmd := &cobra.Command{
		Use:   use,
		Short: short,
		Long:  long,
		CompletionOptions: cobra.CompletionOptions{
			HiddenDefaultCmd: true,
		},
	}

	cmd.SetOut(c.writer)

	return cmd
}
