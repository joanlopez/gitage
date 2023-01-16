package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/joanlopez/gitage/internal/fs"
	"github.com/joanlopez/gitage/internal/log"
)

type CLI struct {
	ctx context.Context
	fs  fs.FS

	// Flags
	path           string
	recipients     []string
	identitiesPath string

	// Writer
	writer log.Writer

	// Commands
	root       *cobra.Command
	init       *cobra.Command
	register   *cobra.Command
	unregister *cobra.Command
	encrypt    *cobra.Command
	decrypt    *cobra.Command
}

func New(ctx context.Context, fs fs.FS) *CLI {
	c := &CLI{
		ctx: ctx,
		fs:  fs,

		// Writer
		writer: log.For(ctx),

		// Flags
		recipients: make([]string, 0),
	}

	c.rootCmd().AddCommand(c.initCmd())
	c.rootCmd().AddCommand(c.registerCmd())
	c.rootCmd().AddCommand(c.unregisterCmd())
	c.rootCmd().AddCommand(c.encryptCmd())
	c.rootCmd().AddCommand(c.decryptCmd())

	return c
}

func (c *CLI) Execute(args ...string) error {
	c.rootCmd().SetArgs(args)
	return c.rootCmd().Execute()
}
