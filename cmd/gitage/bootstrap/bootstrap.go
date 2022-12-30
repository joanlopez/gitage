package bootstrap

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/joanlopez/gitage/internal/fs"
	"github.com/spf13/cobra"
)

var (
	_f   fs.FS
	_out io.Writer
)

func Run(f fs.FS, out io.Writer, args ...string) {
	_f = f
	_out = out
	os.Args = append([]string{"gitage"}, args...)

	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintf(out, "error: %s", err)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&path, "path", "p", "", "path to the repository")
	rootCmd.AddCommand(initCmd)
}

var (
	path string

	rootCmd = &cobra.Command{
		Use:   "gitage",
		Short: "Git+Age = Gitage; simple, modern and secure Git encryption tool",
		Long: `Gitage is a CLI tool that can be used as a wrapper of Git CLI.
It uses 'age' encryption tool to encrypt files before committing them to the repository.`,
		CompletionOptions: cobra.CompletionOptions{
			HiddenDefaultCmd: true,
		},
		Version: "v0.1.0",
	}

	initCmd = &cobra.Command{
		Use:   "init",
		Short: "Initialize a new gitage repository",
		Long: `init is for initializing a new gitage repository.
It creates the .gitage directory with the subsequent files.`,
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			wd, err := os.Getwd()
			if err != nil {
				return err
			}

			if path != "" {
				wd, err = filepath.Abs(path)
				if err != nil {
					return err
				}
			}

			gitageDir := filepath.Join(wd, ".gitage")

			_, _ = fmt.Fprintf(_out, "Creating %s directory...\n", gitageDir)

			if err := fs.Mkdir(_f, gitageDir); err != nil {
				return err
			}

			if err := fs.Create(_f, filepath.Join(gitageDir, "config"), nil); err != nil {
				return err
			}

			if err := fs.Create(_f, filepath.Join(gitageDir, "recipients"), nil); err != nil {
				return err
			}

			_, _ = fmt.Fprintln(_out, "Gitage repository initialized with success!")

			return nil
		},
	}
)
