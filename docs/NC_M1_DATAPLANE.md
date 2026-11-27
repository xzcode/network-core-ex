# NC-M1 data-plane extraction (2026-11-27)

Development branch only. Source: Network EX security candidate 6893c401d137696840a6fb198918eaf29ba5f128; the previous standard-library protocol and discovery import is retained. Changes are namespace relocation, relay statistics types, public SDK wrappers, standalone CLI extraction, and the integration script's executable names. Product account, management, enrollment, daemon, node registration and UI remain product-owned. Wire protocol/signature labels remain unchanged for compatibility; do not share issuer keys between products.

## Interfaces

- pkg/engine: Config, Peer, External, Protector, Status, Check, Load, Run. Run owns the tunnel until cancellation; callers must serialize reconfiguration on a given interface. Managed fields excluded from JSON must come from an authenticated, authorized product adapter.
- pkg/relay: lease issue/verification and authenticated relay/signaling HTTP handlers. Product control planes own room grants and issuer private keys. The relay receives only public issuer keys during service execution.
- pkg/stun: STUN binding listener; pkg/adapter: native VPN packet queue and detached FD adapter.
- cmd/agent: static keygen/check/run/status.
- cmd/relay: explicit TLS listener with /relay/v1 and /signal/v1; optional STUN; issuer-keygen/lease for isolated development or trusted administration. Default TLS listener is loopback; no anonymous relay, no account API.

## Verification

Tokyo Linux/amd64, Go 1.26.11: go mod tidy, go test -count=1 ./..., go vet ./..., and native builds of agent and relay completed. Windows amd64 and macOS arm64 cross-builds of all packages and both final CLI programs passed (not native installation). govulncheck ./... completed with "No vulnerabilities found" at this revision; this is a dependency scan, not a full security audit.
Build the two binaries and run the isolated Linux test:
```sh
mkdir -p bin
go build -o bin/network-core-agent ./cmd/agent
go build -o bin/network-core-relay ./cmd/relay
sudo scripts/m1-netns.sh outside "$(pwd)/bin"
```
Actual final test exit 0: three isolated nodes; simultaneous configured UDP direct and authenticated WSS relay ICMP, 1200-byte TCP/UDP echoes on both paths; application marker absent from 58 captured underlay frames; WireGuard handshake and traffic counters verified; relay restart recovered same overlay IP; unknown WireGuard key rejected; shutdown removed virtual interface and overlay route. Generated ephemeral keys and temporary files cleaned by the test; no public service deployed.

This does not validate complex NAT traversal, Windows TAP, commercial game compatibility, cross-room authorization, SDK stability or production security. The upstream security candidate still has intermittent strict double-NAT application traffic failures. Keep this branch separate until those gates pass and the product integration is verified. Project license remains owner-pending.
