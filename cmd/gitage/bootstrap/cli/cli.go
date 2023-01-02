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
	path           string
	recipients     []string
	identitiesPath string

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

		// Flags
		recipients: make([]string, 0),
	}

	// Root
	c.rootCmd().PersistentFlags().StringVarP(&c.path, "path", "p", "", "path to the repository")
	c.rootCmd().AddCommand(c.initCmd())
	c.rootCmd().AddCommand(c.registerCmd())
	c.rootCmd().AddCommand(c.unregisterCmd())
	c.rootCmd().AddCommand(c.encryptCmd())
	c.rootCmd().AddCommand(c.decryptCmd())

	// Init
	c.initCmd().Flags().StringArrayVarP(&c.recipients, "recipient", "r", nil, "recipients to encrypt the repository")

	// Register
	c.registerCmd().Flags().StringArrayVarP(&c.recipients, "recipient", "r", nil, "recipients to encrypt the repository")
	if err := c.registerCmd().MarkFlagRequired("recipient"); err != nil {
		panic(err)
	}

	// Unregister
	c.unregisterCmd().Flags().StringArrayVarP(&c.recipients, "recipient", "r", nil, "recipients to encrypt the repository")
	if err := c.unregisterCmd().MarkFlagRequired("recipient"); err != nil {
		panic(err)
	}

	// Encrypt
	c.encryptCmd().Flags().StringArrayVarP(&c.recipients, "recipient", "r", nil, "recipients to encrypt the repository")
	if err := c.encryptCmd().MarkFlagRequired("recipient"); err != nil {
		panic(err)
	}

	// Decrypt
	c.decryptCmd().Flags().StringVarP(&c.identitiesPath, "identities", "i", "", "path to the identities file")
	if err := c.decryptCmd().MarkFlagRequired("identities"); err != nil {
		panic(err)
	}

	return c
}

func (c *CLI) Execute(args ...string) error {
	c.rootCmd().SetArgs(args)
	return c.rootCmd().Execute()
}
