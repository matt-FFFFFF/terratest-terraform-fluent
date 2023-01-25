package setuptest

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateAzureRMProvidersFile(t *testing.T) {
	t.Parallel()

	test := Dirs("testdata/depth1", "").WithVars(map[string]interface{}{}).InitAndPlanAndShowWithStruct(t)
	require.NoError(t, test.Err)
	test.CreateAzureRMProvidersFile()
	fp := filepath.Join(test.TmpDir, azureRmfileName)
	require.FileExists(t, fp)
	f, err := os.Open(fp)
	require.NoError(t, err)
	defer f.Close()
	contents, err := io.ReadAll(f)
	require.NoError(t, err)
	require.Equal(t, azureRmContent, string(contents))
}
