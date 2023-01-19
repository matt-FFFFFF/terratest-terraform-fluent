package setuptest

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCopyFilesToTempAndCleanupDepth1(t *testing.T) {
	t.Parallel()

	tmp, cleanup, err := CopyTerraformFolderToTempAndCleanUp(t, "testdata/depth1", "")
	assert.NoError(t, err)
	assert.DirExists(t, tmp)
	cleanup()
	parent := filepath.Dir(tmp)
	t.Logf("parent: %s", parent)
	assert.NoDirExists(t, parent)
}

func TestCopyFilesToTempAndCleanupDepth2(t *testing.T) {
	t.Parallel()

	tmp, cleanup, err := CopyTerraformFolderToTempAndCleanUp(t, "testdata/depth2", "subdir")
	assert.NoError(t, err)
	assert.DirExists(t, tmp)
	cleanup()
	parent := filepath.Dir(filepath.Dir(tmp))
	t.Logf("parent: %s", parent)
	assert.NoDirExists(t, parent)
}

func TestCopyFilesToTempAndCleanupDepth3(t *testing.T) {
	t.Parallel()

	tmp, cleanup, err := CopyTerraformFolderToTempAndCleanUp(t, "testdata/depth3", "subdir/subdir2")
	assert.NoError(t, err)
	assert.DirExists(t, tmp)
	cleanup()
	parent := filepath.Dir(filepath.Dir(filepath.Dir(tmp)))
	t.Logf("parent: %s", parent)
	assert.NoDirExists(t, parent)
}
