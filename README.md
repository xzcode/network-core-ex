# Network Core EX

Shared networking engine for [Network EX](https://github.com/xzcode/network-ex) and [Battle EX](https://github.com/xzcode/battle-ex).

**Status: NC-M1 data-plane candidate / development only.** No production release, stable SDK, or game compatibility claim.

Initial import: [5ebec96](https://github.com/xzcode/network-core-ex/commit/5ebec96758bcc7261a5c17461174898df5b61d08), [migration map](docs/MIGRATION.md). Standard-library-only protocol/policy and discovery codec packages have been extracted; Go 1.26.11 on Linux/amd64 verified `go test -count=1 ./...`, `go vet ./...`, and `go test -race -count=1 ./...` successfully on 2026-11-27 against the equivalent extracted source. WireGuard, ICE, relay, virtual interfaces, product integration and real-game acceptance remain pending.

The additional extracted engine, standalone programs and isolated-network test are documented in [NC-M1 data plane](docs/NC_M1_DATAPLANE.md). This branch is not approved for production or merged into master.

## Architecture

- Go core; reuse WireGuard and Pion rather than create cryptography or NAT traversal from scratch.
- Control plane remains product-owned. Network EX owns device registration, private networks and account policy; Battle EX owns users, lobby, rooms, matchmaking and results.
- Core owns encrypted peer transport, per-peer direct/relay selection, virtual interfaces, discovery and Ethernet packet handling, diagnostics, and reusable authenticated relay/STUN components.
- Product adapters supply authenticated short-lived network membership, peer snapshots, relay grants and discovery policy. Product control planes must not share user identity or trust domains by default.
- Go package SDK first, process/IPC boundary for desktop sidecars and mobile bindings later. No duplicate network engine or mandatory public service.

## Migration baseline

Pinned upstream: [Network EX 792cb70](https://github.com/xzcode/network-ex/commit/792cb70f54fd2c687080a55d36457dfb01320e57).
Existing upstream Linux tests cover IPv4, mixed peer paths, controlled UDP discovery and Linux TAP/ARP. Actual games, Windows/macOS TAP and mobile devices remain unverified.
The separate [M9 security candidate](https://github.com/xzcode/network-ex/tree/codex/m9-security-candidate) is not approved for merge: dual-NAT reliability remains unresolved. Do not ship the old vulnerable dependency baseline or silently replace the current upstream implementation.

## Roadmap

| Milestone | Deliverable |
| --- | --- |
| NC-M0 | Freeze module ownership, dependency boundary and migration gates |
| NC-M1 | Extract protocol, policy and packet codecs first; then encrypted transport, TUN/TAP, ICE and relay; preserve wire behavior and tests |
| NC-M2 | Expose versioned SDK; migrate Network EX on an integration branch, compare existing real-network evidence |
| NC-M3 | Battle discovery profiles, controlled L2/L3 compatibility, game-specific adapter interface |
| NC-M4 | Room-scoped membership, authorization, isolation, revocation and latency diagnostics |
| NC-M5 | Battle EX integration and real Windows game tests |
| NC-M6 | Regional node selection, capacity and fault tolerance |

No changes to Network EX master during extraction. Preserve rollback to pinned upstream and run source-equivalent tests before integration. No anonymous relay, default unrestricted broadcast, or plaintext fallback. Room isolation is enforced by authorized peers and packet policy, not a room ID supplied by the client. License is pending owner confirmation; public source visibility is not an open-source license.

## 中文

Network Core EX 是两个产品共享的 Go 网络引擎。业务控制面与数据面分离，优先迁移已实现模块，不重写密码学与穿透协议。先独立验收再接入 Network EX / Battle EX，真实游戏联机、Windows TAP、双 NAT 及安全门槛必须单独验证。
