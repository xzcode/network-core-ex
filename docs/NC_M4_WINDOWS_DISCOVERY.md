# NC-M4 Windows directed IPv4 discovery candidate

Based on NC-M3 path-refresh candidate 0c5a014ae732e8dd67d1517263d952905dec7843.
This is a development branch, not a stable SDK or Windows game certification.

The existing authenticated UDP discovery codec is reused; no new packet
engine or transparent physical-NIC bridge is introduced. Windows allows only
the dedicated room /24 directed broadcast, with exact UDP port policy
supplied by the trusted product. Limited broadcast 255.255.255.255,
multicast and IPX are not enabled on Windows. The Windows preflight checks
the full /24 for conflicting non-default routes before adding a connected
room subnet. The Wintun address prefix is changed to /24 for that discovery
profile only; WireGuard peer AllowedIPs and the packet guard still restrict
unicast to authorized /32s and the codec restricts discovery to approved
ports. Native stack delivery and application interface selection on Windows
need real-device validation. Linux retains the existing exact discovery routes.

Underlay endpoint protection includes local/peer overlay addresses and the
discovery subnet; the product can protect its HTTPS authority host.
Transport-only refresh may retain an unchanged discovery codec/policy,
virtual adapter and peer keys. Windows route setup/cleanup batches fixed
commands scoped to the owned interface, with a bounded PowerShell execution
time and partial-setup cleanup record. A canceled setup does not bring up
WireGuard.

Evidence: Go unit tests/vet and Windows amd64 cross-build; Battle EX isolated
Linux integration exercised real UDP directed broadcast on allowlisted
port 6112 across both ICE direct and WSS fallback, blocked another port and
a separate room, and completed relay migration, lease refresh with an
existing TCP socket, member revocation and cleanup after authority loss.
No Windows Wintun, Windows route, IPX, L2/TAP or real game certification.
