// Package stun provides the extracted STUN binding service. The caller owns
// the listener and deployment policy.
package stun

import (
	"context"
	"net"

	impl "github.com/xzcode/network-core-ex/internal/node/stun"
)

func Serve(ctx context.Context, conn *net.UDPConn) error { return impl.Serve(ctx, conn) }
