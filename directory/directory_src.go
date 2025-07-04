package directory

import (
	"context"
	"io"
	"os"

	"github.com/loft-sh/image/internal/imagesource/impl"
	"github.com/loft-sh/image/internal/imagesource/stubs"
	"github.com/loft-sh/image/internal/manifest"
	"github.com/loft-sh/image/internal/private"
	"github.com/loft-sh/image/types"
	"github.com/opencontainers/go-digest"
)

type dirImageSource struct {
	impl.Compat
	impl.PropertyMethodsInitialize
	impl.DoesNotAffectLayerInfosForCopy
	stubs.NoGetBlobAtInitialize

	ref dirReference
}

// newImageSource returns an ImageSource reading from an existing directory.
// The caller must call .Close() on the returned ImageSource.
func newImageSource(ref dirReference) private.ImageSource {
	s := &dirImageSource{
		PropertyMethodsInitialize: impl.PropertyMethods(impl.Properties{
			HasThreadSafeGetBlob: false,
		}),
		NoGetBlobAtInitialize: stubs.NoGetBlobAt(ref),

		ref: ref,
	}
	s.Compat = impl.AddCompat(s)
	return s
}

// Reference returns the reference used to set up this source, _as specified by the user_
// (not as the image itself, or its underlying storage, claims).  This can be used e.g. to determine which public keys are trusted for this image.
func (s *dirImageSource) Reference() types.ImageReference {
	return s.ref
}

// Close removes resources associated with an initialized ImageSource, if any.
func (s *dirImageSource) Close() error {
	return nil
}

// GetManifest returns the image's manifest along with its MIME type (which may be empty when it can't be determined but the manifest is available).
// It may use a remote (= slow) service.
// If instanceDigest is not nil, it contains a digest of the specific manifest instance to retrieve (when the primary manifest is a manifest list);
// this never happens if the primary manifest is not a manifest list (e.g. if the source never returns manifest lists).
func (s *dirImageSource) GetManifest(ctx context.Context, instanceDigest *digest.Digest) ([]byte, string, error) {
	path, err := s.ref.manifestPath(instanceDigest)
	if err != nil {
		return nil, "", err
	}
	m, err := os.ReadFile(path)
	if err != nil {
		return nil, "", err
	}
	return m, manifest.GuessMIMEType(m), err
}

// GetBlob returns a stream for the specified blob, and the blob’s size (or -1 if unknown).
// The Digest field in BlobInfo is guaranteed to be provided, Size may be -1 and MediaType may be optionally provided.
// May update BlobInfoCache, preferably after it knows for certain that a blob truly exists at a specific location.
func (s *dirImageSource) GetBlob(ctx context.Context, info types.BlobInfo, cache types.BlobInfoCache) (io.ReadCloser, int64, error) {
	path, err := s.ref.layerPath(info.Digest)
	if err != nil {
		return nil, -1, err
	}
	r, err := os.Open(path)
	if err != nil {
		return nil, -1, err
	}
	fi, err := r.Stat()
	if err != nil {
		return nil, -1, err
	}
	return r, fi.Size(), nil
}
