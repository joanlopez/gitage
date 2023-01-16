package bootstrap_test

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"filippo.io/age"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/joanlopez/gitage"
	"github.com/joanlopez/gitage/cmd/gitage/bootstrap"
	"github.com/joanlopez/gitage/internal/fs"
	"github.com/joanlopez/gitage/internal/fs/archive"
	"github.com/joanlopez/gitage/internal/log"
)

func Test(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		dir  string
		args []string
	}{
		// ~/$ gitage
		{dir: "no-cmd-no-args", args: []string{}},

		// ~/$ gitage init
		{dir: "init-empty-repo", args: []string{"init", "-p", "/repo"}},
		{dir: "init-existing-repo", args: []string{"init", "-p", "/repo"}},
		{dir: "init-wrong-repo", args: []string{"init", "-p", "/repo"}},
		{dir: "init-repo-with-single-recipient", args: []string{"init", "-p", "/repo", "-r", "age1ql3z7hjy54pw3hyww5ayyfg7zqgvc7w3j2elw8zmrj2kg5sfn9aqmcac8p"}},
		{dir: "init-repo-with-multiple-recipients", args: []string{"init", "-p", "/repo", "-r", "age1ql3z7hjy54pw3hyww5ayyfg7zqgvc7w3j2elw8zmrj2kg5sfn9aqmcac8p", "-r", "age1lggyhqrw2nlhcxprm67z43rta597azn8gknawjehu9d9dl0jq3yqqvfafg"}},

		// ~/$ gitage register
		{dir: "register-no-args", args: []string{"register"}},
		{dir: "register-no-recipients", args: []string{"register", "-p", "/repo"}},
		{dir: "register-empty-repo", args: []string{"register", "-p", "/repo", "-r", "age1ql3z7hjy54pw3hyww5ayyfg7zqgvc7w3j2elw8zmrj2kg5sfn9aqmcac8p"}},
		{dir: "register-first-recipient", args: []string{"register", "-p", "/repo", "-r", "age1ql3z7hjy54pw3hyww5ayyfg7zqgvc7w3j2elw8zmrj2kg5sfn9aqmcac8p"}},
		{dir: "register-repeated-recipient", args: []string{"register", "-p", "/repo", "-r", "age1ql3z7hjy54pw3hyww5ayyfg7zqgvc7w3j2elw8zmrj2kg5sfn9aqmcac8p"}},
		{dir: "register-single-recipient", args: []string{"register", "-p", "/repo", "-r", "age1lggyhqrw2nlhcxprm67z43rta597azn8gknawjehu9d9dl0jq3yqqvfafg"}},
		{dir: "register-multiple-recipients", args: []string{"register", "-p", "/repo", "-r", "age1lggyhqrw2nlhcxprm67z43rta597azn8gknawjehu9d9dl0jq3yqqvfafg", "-r", "age1yhm4gctwfmrpz87tdslm550wrx6m79y9f2hdzt0lndjnehwj0ukqrjpyx5"}},

		// ~/$ gitage unregister
		{dir: "unregister-no-args", args: []string{"unregister"}},
		{dir: "unregister-no-recipients", args: []string{"unregister", "-p", "/repo"}},
		{dir: "unregister-empty-repo", args: []string{"unregister", "-p", "/repo", "-r", "age1ql3z7hjy54pw3hyww5ayyfg7zqgvc7w3j2elw8zmrj2kg5sfn9aqmcac8p"}},
		{dir: "unregister-single-recipient", args: []string{"unregister", "-p", "/repo", "-r", "age1lggyhqrw2nlhcxprm67z43rta597azn8gknawjehu9d9dl0jq3yqqvfafg"}},
		{dir: "unregister-multiple-recipients", args: []string{"unregister", "-p", "/repo", "-r", "age1lggyhqrw2nlhcxprm67z43rta597azn8gknawjehu9d9dl0jq3yqqvfafg", "-r", "age1yhm4gctwfmrpz87tdslm550wrx6m79y9f2hdzt0lndjnehwj0ukqrjpyx5"}},
		{dir: "unregister-last-recipient", args: []string{"unregister", "-p", "/repo", "-r", "age1ql3z7hjy54pw3hyww5ayyfg7zqgvc7w3j2elw8zmrj2kg5sfn9aqmcac8p"}},
		{dir: "unregister-last-recipient", args: []string{"unregister", "-p", "/repo", "-r", "age1ql3z7hjy54pw3hyww5ayyfg7zqgvc7w3j2elw8zmrj2kg5sfn9aqmcac8p"}},

		// ~/$ gitage encrypt
		{dir: "encrypt-multiple-files", args: []string{"encrypt", "-p", "/repo/data", "-r", "age1xkt49yr0y689x45qqrja6rgl0sne82gw5gt6mhhepa7xm7r6myfsd63983"}},

		// ~/$ gitage decrypt
		{dir: "decrypt-multiple-files", args: []string{"decrypt", "-p", "/repo/data", "-i", "/repo/.gitage/identities"}},
	}

	for _, tc := range tcs {
		tc := tc

		t.Run(tc.dir, func(t *testing.T) {
			// Different test cases can be executed in parallel
			t.Parallel()

			// Create a new filesystem
			f := fsForTestCase(t, tc.dir)

			// Create a new buffer to capture the output
			out := new(bytes.Buffer)
			ctx := log.Ctx(out)

			// Run the bootstrap
			bootstrap.Run(ctx, f, tc.args...)

			// Assert the results
			ass := newAsserter(t, tc.dir, f, out)
			ass.assertOutput()
			ass.assertFileTree()
		})
	}
}

func fsForTestCase(t *testing.T, dirName string) fs.FS {
	t.Helper()

	const initFilePathFmt = "./testdata/%s/init.txtar"
	initFilePath, err := filepath.Abs(fmt.Sprintf(initFilePathFmt, dirName))
	require.NoError(t, err)

	info, err := os.Stat(initFilePath)
	if err == nil && !info.IsDir() {
		const filename = "init"
		return fsFromTxtarFile(t, dirName, filename)
	}

	memFS := afero.NewMemMapFs()

	const initDirPathFmt = "./testdata/%s/init"
	initDirPath, err := filepath.Abs(fmt.Sprintf(initDirPathFmt, dirName))
	require.NoError(t, err)

	err = filepath.Walk(initDirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath := relPath(t, initDirPath, path)

		if info.IsDir() {
			return fs.Mkdir(memFS, relPath)
		}

		contents, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		if err = fs.Create(memFS, relPath, contents); err != nil {
			return err
		}

		return nil
	})

	require.NoError(t, err)

	return memFS
}

func fsFromTxtarFile(t *testing.T, dir, filename string) fs.FS {
	t.Helper()

	const initFilePathFmt = "./testdata/%s/%s.txtar"
	initFilePath, err := filepath.Abs(fmt.Sprintf(initFilePathFmt, dir, filename))
	require.NoError(t, err)

	b, err := os.ReadFile(initFilePath)
	require.NoError(t, err)

	f, err := fs.FromArchive(archive.Parse(b))
	require.NoError(t, err)

	return f
}

type asserter struct {
	t   *testing.T
	dir string

	testArchive *archive.Archive
	testOut     string

	expectedArchive *archive.Archive
	expectedOut     string

	identities []age.Identity
}

func newAsserter(t *testing.T, dir string, testFS fs.FS, testOut *bytes.Buffer) asserter {
	t.Helper()

	testArchive, err := fs.ToArchive(testFS)
	require.NoError(t, err)

	expectedArchive, err := fs.ToArchive(fsFromTxtarFile(t, dir, "expected"))
	require.NoError(t, err)

	const outputFilePathFmt = "./testdata/%s/out.txt"
	outFilePath, err := filepath.Abs(fmt.Sprintf(outputFilePathFmt, dir))
	require.NoError(t, err)

	b, err := os.ReadFile(outFilePath)
	require.NoError(t, err)

	return asserter{
		t:   t,
		dir: dir,

		testArchive: testArchive,
		testOut:     testOut.String(),

		expectedArchive: expectedArchive,
		expectedOut:     string(b),

		identities: identitiesFromFile(t, dir),
	}
}

func (a asserter) assertFileTree() {
	a.t.Helper()

	// We do use a map of booleans to check which paths
	// have been visited (ergo, are expected) and which
	// have not (ergo, are unexpected).
	visited := make(map[string]bool)

	// All expected files should be present in the got filesystem.
	// Either encrypted (in which case we compare the decrypted content)
	// or not encrypted (in which case we compare the content as is).
	for f := range a.expectedArchive.Files() {
		visited[f.Name] = true

		var gotData string
		switch filepath.Ext(f.Name) {
		case gitage.Ext:
			decrypted, err := gitage.Decrypt(context.Background(), fileContents(a.t, a.testArchive, f), a.identities...)
			assert.NoError(a.t, err, "Failed to decrypt file from test file system: %s", f.Name)
			gotData = string(decrypted)
		default:
			gotData = string(fileContents(a.t, a.testArchive, f))
		}

		assert.Equal(a.t, string(f.Data), gotData, "File content was not as expected: %s", f.Name)
	}

	// All got files should have been visited already (expected).
	// Otherwise, we have a file that was not expected.
	for f := range a.testArchive.Files() {
		_, visited := visited[f.Name]
		assert.True(a.t, visited, "File from test file system was not expected: %s", f.Name)
	}
}

func (a asserter) assertOutput() {
	a.t.Helper()
	assert.Equal(a.t, a.expectedOut, a.testOut, "Execution output was not as expected")
}

func identitiesFromFile(t *testing.T, dir string) []age.Identity {
	t.Helper()

	const identitiesFilePathFmt = "./testdata/%s/identities"
	identitiesFilePath, err := filepath.Abs(fmt.Sprintf(identitiesFilePathFmt, dir))
	require.NoError(t, err)

	f, err := os.Open(identitiesFilePath)
	if err != nil && os.IsNotExist(err) {
		return nil
	}
	require.NoError(t, err)

	identities, err := age.ParseIdentities(f)
	require.NoError(t, err)

	return identities
}

func fileContents(t *testing.T, archive *archive.Archive, f *archive.File) []byte {
	t.Helper()
	file := archive.Get(f.Name)
	require.NotNil(t, file, "File not found in archive: %s", f.Name)
	return file.Data
}

func relPath(t *testing.T, root, path string) string {
	t.Helper()

	path = strings.Replace(path, root, "", 1)
	if path == "" {
		return "/"
	}

	return path
}
