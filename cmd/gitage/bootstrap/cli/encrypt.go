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
		c.encrypt = &cobra.Command{
			Use:   "encrypt",
			Short: "Encrypts files on the specified path",
			Args:  cobra.ExactArgs(0),
			RunE: func(cmd *cobra.Command, args []string) error {
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
			},
		}
	}

	return c.encrypt
}
