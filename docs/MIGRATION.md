# NC-M1 extraction: first slice (2026-11-27)

Upstream baseline: [792cb70](https://github.com/xzcode/network-ex/commit/792cb70f54fd2c687080a55d36457dfb01320e57). This is a development import, not the stable network SDK.

| Upstream location | Core location | Scope |
|---|---|---|
| pkg/protocol/{advanced,discovery,ethernet,exit}.go and corresponding tests | pkg/protocol | Network/service policy, discovery validation, Ethernet MAC/prefix, exit policy |
| internal/dataplane/discovery/codec.go and codec_test.go | pkg/discovery | UDP discovery encapsulation, peer/policy checks and rate limiting |

Only the import path and module layout are changed. This slice uses the Go standard library only; it does not adopt the old ICE/WireGuard dependencies. Existing protocol restrictions (IPv4 subnet sizes, discovery rule limits and MTU) are inherited, not new product limits or proven final game compatibility.

Verification commands from repository root: `go test -count=1 ./...`, `go vet ./...`, `go test -race -count=1 ./...`. Actual execution results are tracked separately; these commands alone are not evidence of real network connectivity.

Remaining NC-M1: separate product-managed types from reusable peer snapshots; extract guard and L2 codec, virtual interfaces, WireGuard transport, authenticated relay and ICE path manager; adapt imports and preserve existing tests. Coordinate dependency selection with the upstream security candidate. Do not treat its unresolved dual-NAT results as passed. Windows adapter, real-game sessions and room isolation still require acceptance.

No modifications to Network EX or Battle EX's runtime. No release, deployment, new project license, or production security claim.
