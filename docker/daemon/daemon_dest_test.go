package daemon

import "github.com/loft-sh/image/internal/private"

var _ private.ImageDestination = (*daemonImageDestination)(nil)
