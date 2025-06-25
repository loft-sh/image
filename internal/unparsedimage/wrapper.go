package unparsedimage

import (
	"github.com/loft-sh/image/internal/private"
	"github.com/loft-sh/image/types"
)

// wrapped provides the private.UnparsedImage operations
// for an object that only implements types.UnparsedImage
type wrapped struct {
	types.UnparsedImage
}

// FromPublic(unparsed) returns an object that provides the private.UnparsedImage API
func FromPublic(unparsed types.UnparsedImage) private.UnparsedImage {
	if unparsed2, ok := unparsed.(private.UnparsedImage); ok {
		return unparsed2
	}
	return &wrapped{
		UnparsedImage: unparsed,
	}
}
