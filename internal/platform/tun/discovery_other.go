//go:build !linux && !windows

package tun

import (
	"errors"
	"github.com/xzcode/network-core-ex/pkg/protocol"
	"net/netip"
)

func ConfigureDiscovery(name string, local netip.Addr, rules []protocol.DiscoveryRule) error {
	if len(rules) > 0 {
		return errors.New("game discovery is currently Linux-only")
	}
	return nil
}

func PreflightDiscovery(name string, local netip.Addr, rules []protocol.DiscoveryRule) error {
	if len(rules) > 0 {
		return errors.New("discovery unsupported on this platform")
	}
	return nil
}

func ValidateDiscovery(local netip.Addr, rules []protocol.DiscoveryRule) error { return nil }
