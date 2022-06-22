package archive

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"code-intelligence.com/cifuzz/util/fileutil"
)

func TestArchive_Create(t *testing.T) {
	// prepare dir to create archive from
	tmpDir, err := os.MkdirTemp("", "cifuzz-archive-test")
	require.NoError(t, err)
	defer fileutil.Cleanup(tmpDir)

	dir := filepath.Join(tmpDir, "foo", "bar")
	err = os.MkdirAll(dir, 0755)
	require.NoError(t, err)

	filename := "hello"
	testData := []byte("foo")
	err = os.WriteFile(filepath.Join(dir, filename), testData, 0644)
	require.NoError(t, err)

	archiveDir, err := os.MkdirTemp("", "cifuzz-archive-test")
	require.NoError(t, err)
	defer fileutil.Cleanup(archiveDir)
	archiveFilepath := filepath.Join(archiveDir, "test.tar.gz")

	err = Create(tmpDir, archiveFilepath)
	require.NoError(t, err)
	require.FileExists(t, archiveFilepath)

	extractDir, err := os.MkdirTemp("", "cifuzz-archive-test")
	require.NoError(t, err)
	defer fileutil.Cleanup(extractDir)

	err = Extract(archiveFilepath, extractDir)
	require.NoError(t, err)

	extractedFile := filepath.Join(extractDir, "foo", "bar", filename)
	require.FileExists(t, extractedFile)
	content, err := os.ReadFile(extractedFile)
	require.NoError(t, err)
	assert.Equal(t, testData, content)
}

// Creates the archive in the source directory
func TestArchive_CreateInSourceDir(t *testing.T) {
	// prepare dir to create archive from
	tmpDir, err := os.MkdirTemp("", "cifuzz-archive-test")
	require.NoError(t, err)
	defer fileutil.Cleanup(tmpDir)

	filename := "test.txt"
	testData := []byte("foo")
	err = os.WriteFile(filepath.Join(tmpDir, filename), testData, 0644)
	require.NoError(t, err)

	archiveFilepath := filepath.Join(tmpDir, "test.tar.gz")

	// create the archive in the same directory
	err = Create(tmpDir, archiveFilepath)
	require.NoError(t, err)
	require.FileExists(t, archiveFilepath)

	extractDir, err := os.MkdirTemp("", "cifuzz-archive-test")
	require.NoError(t, err)
	defer fileutil.Cleanup(extractDir)

	err = Extract(archiveFilepath, extractDir)
	require.NoError(t, err)

	extractedFile := filepath.Join(extractDir, filename)
	require.FileExists(t, extractedFile)
	content, err := os.ReadFile(extractedFile)
	require.NoError(t, err)
	assert.Equal(t, testData, content)
}

// Extract an existing test archive
func TestArchive_Extract(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cifuzz-archive-test")
	require.NoError(t, err)
	defer fileutil.Cleanup(tmpDir)

	err = Extract("testdata/test.tar.gz", tmpDir)
	require.NoError(t, err)

	helloFile := filepath.Join(tmpDir, "hello.txt")
	require.FileExists(t, helloFile)
	content, err := os.ReadFile(helloFile)
	require.NoError(t, err)
	assert.Equal(t, "world\n", string(content))

	require.DirExists(t, filepath.Join(tmpDir, "foo"))
	require.FileExists(t, filepath.Join(tmpDir, "foo", "bar"))
}
