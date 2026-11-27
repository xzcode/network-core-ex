// Package adapter exposes native VPN packet interfaces for product-owned
// mobile and desktop integrations.
package adapter

import (
	"github.com/xzcode/network-core-ex/internal/platform/packetfd"
	"github.com/xzcode/network-core-ex/internal/platform/packetqueue"
	"golang.zx2c4.com/wireguard/tun"
)

type PacketQueue = packetqueue.Device

func NewPacketQueue(write func([]byte) bool, closeNative func()) *PacketQueue {
	return packetqueue.New(write, closeNative)
}

// OpenPacketFD takes ownership of a valid detached native VPN descriptor on Unix, even if initialization fails.
func OpenPacketFD(fd int, name string) (tun.Device, error) { return packetfd.New(fd, name) }
