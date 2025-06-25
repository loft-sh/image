//go:build containers_image_storage_stub

package alltransports

import "github.com/loft-sh/image/transports"

func init() {
	transports.Register(transports.NewStubTransport("containers-storage"))
}
