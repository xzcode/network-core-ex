//go:build windows

package tun

import (
	"errors"
	"github.com/xzcode/network-core-ex/pkg/protocol"
	"net/netip"
	"strconv"
)

// Windows Wintun is an IPv4 L3 adapter, not a transparent Ethernet bridge.
// Only the room's directed /24 broadcast is eligible. Never redirect
// 255.255.255.255, multicast, IPX, or the physical LAN.
func ValidateDiscovery(local netip.Addr, rules []protocol.DiscoveryRule) error {
	if len(rules) == 0 {
		return nil
	}
	if !local.Is4() {
		return errors.New("IPv4 discovery required")
	}
	b := local.As4()
	b[3] = 255
	broadcast := netip.AddrFrom4(b)
	for _, r := range rules {
		if r.Destination != broadcast.String() || r.Port == 0 {
			return errors.New("Windows discovery requires the room directed broadcast")
		}
	}
	return nil
}
func roomTargets(local netip.Addr) []netip.Addr {
	b := local.As4()
	targets := make([]netip.Addr, 0, 256)
	for i := 0; i < 256; i++ {
		b[3] = byte(i)
		targets = append(targets, netip.AddrFrom4(b))
	}
	return targets
}
func PreflightDiscovery(name string, local netip.Addr, rules []protocol.DiscoveryRule) error {
	if len(rules) == 0 {
		return nil
	}
	if e := ValidateDiscovery(local, rules); e != nil {
		return e
	}
	return winOverlap(0, roomTargets(local))
}
func ConfigureDiscovery(name string, local netip.Addr, rules []protocol.DiscoveryRule) error {
	if len(rules) == 0 {
		return nil
	}
	if e := ValidateDiscovery(local, rules); e != nil {
		return e
	}
	winOwned.Lock()
	owned, ok := winOwned.m[name]
	winOwned.Unlock()
	if !ok || owned.address != local.String() {
		return errors.New("discovery interface not owned by engine")
	}
	if e := winOverlap(owned.index, roomTargets(local)); e != nil {
		return e
	}
	// A /24 address tells the Windows IP stack which directed broadcast is
	// local. Its connected route covers ONLY this dedicated room subnet. The
	// existing WireGuard/guard peers still enforce exact /32 authorization.
	_, e := ps("Set-NetIPAddress -InterfaceIndex " + strconv.Itoa(owned.index) + " -IPAddress '" + local.String() + "' -PrefixLength 24 -AddressFamily IPv4 -PolicyStore ActiveStore")
	return e
}
