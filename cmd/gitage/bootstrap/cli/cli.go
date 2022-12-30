package cli

import (
	"context"
	"github.com/spf13/cobra"

	"github.com/joanlopez/gitage/internal/fs"
)

type CLI struct {
	ctx context.Context
	fs  fs.FS

	// Flags
	path       string
	recipients []string

	// Commands
	root       *cobra.Command
	init       *cobra.Command
	register   *cobra.Command
	unregister *cobra.Command
}

func New(ctx context.Context, fs fs.FS) *CLI {
	c := &CLI{
		ctx: ctx,
		fs:  fs,

		// Flags
		recipients: make([]string, 0),
	}

	// Root
	c.rootCmd().PersistentFlags().StringVarP(&c.path, "path", "p", "", "path to the repository")
	c.rootCmd().AddCommand(c.initCmd())
	c.rootCmd().AddCommand(c.registerCmd())
	c.rootCmd().AddCommand(c.unregisterCmd())

	// Init
	c.initCmd().Flags().StringArrayVarP(&c.recipients, "recipient", "r", nil, "recipients to encrypt the repository")

	// Register
	c.registerCmd().Flags().StringArrayVarP(&c.recipients, "recipient", "r", nil, "recipients to encrypt the repository")
	c.registerCmd().MarkFlagRequired("recipient")

	// Unregister
	c.unregisterCmd().Flags().StringArrayVarP(&c.recipients, "recipient", "r", nil, "recipients to encrypt the repository")
	c.unregisterCmd().MarkFlagRequired("recipient")

	return c
}

func (c *CLI) Execute(args ...string) error {
	c.rootCmd().SetArgs(args)
	return c.rootCmd().Execute()
}
