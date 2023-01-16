package cli

import "github.com/spf13/cobra"

func (c *CLI) rootCmd() *cobra.Command {
	if c.root == nil {
		// Init base command
		c.root = c.command(
			"gitage",
			"Git+age = Gitage; simple, modern and secure Git encryption tool",
			`Gitage is a CLI tool that can be used as a wrapper of Git CLI.
It uses 'age' encryption tool to encrypt files before committing them to the repository.`,
		)

		// Set version
		c.root.Version = "v0.1.0"

		// Set flags
		c.root.PersistentFlags().StringVarP(&c.path, "path", "p", "", "path to the repository")
	}

	return c.root
}
