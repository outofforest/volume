package libnet

import (
	"context"
	"net"
	"time"

	"github.com/ridge/must"
)

var lc = net.ListenConfig{
	KeepAlive: 3 * time.Minute,
}

// Listen creates listener listening on TCP port
func Listen(address string) (net.Listener, error) {
	return lc.Listen(context.Background(), "tcp", address)
}

// ListenOnRandomPort listens on random localhost TCP port
func ListenOnRandomPort() net.Listener {
	return must.NetListener(Listen("localhost:"))
}
