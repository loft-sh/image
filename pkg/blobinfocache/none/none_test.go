package none

import (
	"github.com/loft-sh/image/types"
)

var _ types.BlobInfoCache = &noCache{}
