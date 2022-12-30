package bootstrap_test

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/txtar"

	"github.com/joanlopez/gitage/cmd/gitage/bootstrap"
	"github.com/joanlopez/gitage/internal/fs"
	"github.com/joanlopez/gitage/internal/log"
)

func Test(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		dir  string
		args []string
	}{
		{dir: "init-empty-repo", args: []string{"init", "-p", "/repo"}},
		{dir: "init-existing-repo", args: []string{"init", "-p", "/repo"}},
		{dir: "init-wrong-repo", args: []string{"init", "-p", "/repo"}},
		{dir: "init-repo-with-single-recipient", args: []string{"init", "-p", "/repo", "-r", "age1ql3z7hjy54pw3hyww5ayyfg7zqgvc7w3j2elw8zmrj2kg5sfn9aqmcac8p"}},
		{dir: "init-repo-with-multiple-recipients", args: []string{"init", "-p", "/repo", "-r", "age1ql3z7hjy54pw3hyww5ayyfg7zqgvc7w3j2elw8zmrj2kg5sfn9aqmcac8p", "-r", "age1lggyhqrw2nlhcxprm67z43rta597azn8gknawjehu9d9dl0jq3yqqvfafg"}},
		// register no args
		{dir: "register-empty-repo", args: []string{"register", "-p", "/repo", "-r", "age1ql3z7hjy54pw3hyww5ayyfg7zqgvc7w3j2elw8zmrj2kg5sfn9aqmcac8p"}},
		{dir: "register-single-recipient", args: []string{"register", "-p", "/repo", "-r", "age1lggyhqrw2nlhcxprm67z43rta597azn8gknawjehu9d9dl0jq3yqqvfafg"}},
		{dir: "register-multiple-recipients", args: []string{"register", "-p", "/repo", "-r", "age1lggyhqrw2nlhcxprm67z43rta597azn8gknawjehu9d9dl0jq3yqqvfafg", "-r", "age1yhm4gctwfmrpz87tdslm550wrx6m79y9f2hdzt0lndjnehwj0ukqrjpyx5"}},
		// unregister no args
		{dir: "unregister-empty-repo", args: []string{"unregister", "-p", "/repo", "-r", "age1ql3z7hjy54pw3hyww5ayyfg7zqgvc7w3j2elw8zmrj2kg5sfn9aqmcac8p"}},
		// unregister single
		// unregister multiple
	}

	for _, tc := range tcs {
		tc := tc

		t.Run(tc.dir, func(t *testing.T) {
			// Different test cases can be executed in parallel
			t.Parallel()

			// Create a new filesystem
			f := fsFromTxtarFile(t, tc.dir, "init")

			// Create a new buffer to capture the output
			out := new(bytes.Buffer)
			ctx := log.Ctx(out)

			// Run the bootstrap
			bootstrap.Run(ctx, f, tc.args...)

			// Expected result
			expected := fsFromTxtarFile(t, tc.dir, "expected")

			// Compare the filesystems
			compareEqualFs(t, expected, f)

			// Compare the output
			compareEqualOutput(t, tc.dir, out)
		})
	}
}

func fsFromTxtarFile(t *testing.T, dir, file string) fs.FS {
	t.Helper()

	b, err := os.ReadFile(fmt.Sprintf("./testdata/%s/%s.txtar", dir, file))
	require.NoError(t, err)

	f, err := fs.FromTxtar(txtar.Parse(b))
	require.NoError(t, err)

	return f
}

func expectedOut(t *testing.T, dir string) []byte {
	t.Helper()

	b, err := os.ReadFile(fmt.Sprintf("./testdata/%s/out.txt", dir))
	require.NoError(t, err)

	return b
}

func compareEqualFs(t *testing.T, fs1, fs2 fs.FS) {
	t.Helper()

	fs1Txtar, err := fs.ToTxtar(fs1)
	require.NoError(t, err)

	fs2Txtar, err := fs.ToTxtar(fs2)
	require.NoError(t, err)

	fs1Bytes := txtar.Format(fs1Txtar)
	fs2Bytes := txtar.Format(fs2Txtar)

	assert.Equal(t, string(fs1Bytes), string(fs2Bytes))
}

func compareEqualOutput(t *testing.T, dir string, out *bytes.Buffer) {
	assert.Equal(t, string(expectedOut(t, dir)), string(out.Bytes()))
}
