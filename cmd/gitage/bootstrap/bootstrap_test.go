package bootstrap

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/joanlopez/gitage/internal/fs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/txtar"
)

func Test(t *testing.T) {
	tcs := []struct {
		dir  string
		args []string
	}{
		{dir: "init-empty-repo", args: []string{"init", "-p", "/repo"}},
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

			// Run the bootstrap
			Run(f, out, tc.args...)

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
