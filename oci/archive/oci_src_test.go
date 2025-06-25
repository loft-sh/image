package archive

import (
	"path/filepath"
	"testing"

	"github.com/loft-sh/image/internal/private"
	"github.com/loft-sh/image/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var _ private.ImageSource = (*ociArchiveImageSource)(nil)

func TestNewImageSourceNotFound(t *testing.T) {
	sysctx := types.SystemContext{}
	emptyDir := t.TempDir()
	archivePath := filepath.Join(emptyDir, "foo.ociarchive")
	imgref, err := ParseReference(archivePath)
	require.NoError(t, err)
	_, err = LoadManifestDescriptorWithContext(&sysctx, imgref)
	assert.NotNil(t, err)
	var aerr ArchiveFileNotFoundError
	assert.ErrorAs(t, err, &aerr)
	assert.Equal(t, aerr.path, archivePath)
}
