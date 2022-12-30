package cli

import "github.com/spf13/cobra"

func (c *CLI) rootCmd() *cobra.Command {
	if c.root == nil {
		c.root = &cobra.Command{
			Use:   "gitage",
			Short: "Git+age = Gitage; simple, modern and secure Git encryption tool",
			Long: `Gitage is a CLI tool that can be used as a wrapper of Git CLI.
It uses 'age' encryption tool to encrypt files before committing them to the repository.`,
			CompletionOptions: cobra.CompletionOptions{
				HiddenDefaultCmd: true,
			},
			Version: "v0.1.0",
		}
	}

	return c.root
}
