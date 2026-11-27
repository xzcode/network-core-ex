// Package engine exposes the provisional shared encrypted network engine.
// Product code must authorize membership before constructing a Config.
package engine

import (
	"context"

	"github.com/xzcode/network-core-ex/internal/dataplane/ethernet"
	"github.com/xzcode/network-core-ex/internal/dataplane/path"
	"github.com/xzcode/network-core-ex/internal/dataplane/tunnel"
	"github.com/xzcode/network-core-ex/internal/platform/protect"
)

type Config = tunnel.Config
type PathUpdate = tunnel.PathUpdate
type Peer = tunnel.Peer
type External = tunnel.External
type Protector = protect.Protector
type Status = tunnel.Status
type PeerStatus = tunnel.PeerStatus
type PathSnapshot = path.Snapshot
type EthernetStats = ethernet.Stats

// Load reads a static configuration. Managed membership and policy must be set
// by the trusted product adapter; fields excluded from JSON are not inferred.
func Load(file string) (Config, error) { return tunnel.Load(file) }

// Check validates the proposed config and checks the relevant key files.
func Check(config Config) error { return tunnel.Check(config) }

// Run owns the tunnel lifecycle until ctx is cancelled or an error occurs.
// Do not run overlapping configurations on the same virtual interface.
func Run(ctx context.Context, config Config) error { return tunnel.Run(ctx, config) }

// RunWithUpdates accepts serialized, trusted transport-only path updates.
// See PathUpdate: the peer membership and packet policy must remain unchanged.
// Failures stop the tunnel; the caller must not retain expired authorization.
func RunWithUpdates(ctx context.Context, config Config, updates <-chan PathUpdate) error {
	return tunnel.RunWithUpdates(ctx, config, updates)
}
