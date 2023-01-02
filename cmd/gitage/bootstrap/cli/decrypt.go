package cli

import (
	"bytes"

	"filippo.io/age"
	"github.com/spf13/cobra"

	"github.com/joanlopez/gitage"
	"github.com/joanlopez/gitage/internal/fs"
	"github.com/joanlopez/gitage/internal/log"
)

func (c *CLI) decryptCmd() *cobra.Command {
	if c.decrypt == nil {
		c.decrypt = &cobra.Command{
			Use:   "decrypt",
			Short: "Decrypts files on the specified path",
			Args:  cobra.ExactArgs(0),
			RunE: func(cmd *cobra.Command, args []string) error {
				rawIdentities, err := fs.Read(c.fs, c.identitiesPath)
				if err != nil {
					return err
				}

				identities, err := age.ParseIdentities(bytes.NewReader(rawIdentities))
				if err != nil {
					return err
				}

				log.For(c.ctx).Println("Decrypting files...")
				err = gitage.DecryptAll(c.ctx, c.fs, c.path, identities...)
				if err != nil {
					return err
				}

				log.For(c.ctx).Println("Files decrypted with success!")

				return err
			},
		}
	}

	return c.decrypt
}
