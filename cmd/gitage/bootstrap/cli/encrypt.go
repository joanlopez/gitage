package cli

import (
	"bytes"
	"strings"

	"filippo.io/age"
	"github.com/spf13/cobra"

	"github.com/joanlopez/gitage"
	"github.com/joanlopez/gitage/internal/log"
)

func (c *CLI) encryptCmd() *cobra.Command {
	if c.encrypt == nil {
		c.encrypt = c.command(
			"encrypt",
			"Encrypts files on the specified path",
			"",
		)

		// Set args
		c.encrypt.Args = cobra.ExactArgs(0)

		// Set flags
		c.encrypt.Flags().StringArrayVarP(&c.recipients, "recipient", "r", nil, "recipients to encrypt the repository")
		if err := c.encrypt.MarkFlagRequired("recipient"); err != nil {
			panic(err)
		}

		// Set run fn
		c.encrypt.RunE = func(cmd *cobra.Command, args []string) error {
			rawRecipients := []byte(strings.Join(c.recipients, "\n"))
			recipients, err := age.ParseRecipients(bytes.NewReader(rawRecipients))
			if err != nil {
				return err
			}

			log.For(c.ctx).Println("Encrypting files...")
			err = gitage.EncryptAll(c.ctx, c.fs, c.path, recipients...)
			if err != nil {
				return err
			}

			log.For(c.ctx).Println("Files encrypted with success!")

			return err
		}
	}

	return c.encrypt
}
