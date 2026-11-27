// Package relay exposes authenticated encrypted-packet relay and signaling.
// The application owns issuer keys, room admission and transport TLS.
package relay

import (
	"crypto/ed25519"
	"time"

	impl "github.com/xzcode/network-core-ex/internal/dataplane/relay"
)

type Lease = impl.Lease
type Server = impl.Server
type SignalServer = impl.SignalServer

func Issue(key ed25519.PrivateKey, lease Lease) (string, error) { return impl.Issue(key, lease) }
func Verify(key ed25519.PublicKey, token string, now time.Time) (Lease, error) {
	return impl.Verify(key, token, now)
}
func NewServer(node string, issuer ed25519.PublicKey) *Server { return impl.NewServer(node, issuer) }
func NewSignalServer(node string, issuer ed25519.PublicKey) *SignalServer {
	return impl.NewSignalServer(node, issuer)
}
