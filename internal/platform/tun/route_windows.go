//go:build windows

package tun

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/netip"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/sys/windows"
)

type windowsRoute struct {
	DestinationPrefix string
	InterfaceIndex    int
}
type ownedWin struct {
	index   int
	address string
	routes  []netip.Addr
}

var winOwned = struct {
	sync.Mutex
	m map[string]ownedWin
}{m: make(map[string]ownedWin)}

func ps(script string) ([]byte, error) {
	sys, e := windows.GetSystemDirectory()
	if e != nil {
		return nil, e
	}
	exe := filepath.Join(sys, "WindowsPowerShell", "v1.0", "powershell.exe")
	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
	defer cancel()
	out, e := exec.CommandContext(ctx, exe, "-NoLogo", "-NoProfile", "-NonInteractive", "-Command", "$ErrorActionPreference='Stop'; "+script).CombinedOutput()
	if e != nil {
		return nil, fmt.Errorf("Windows IP configuration: %w (%s)", e, strings.TrimSpace(string(out)))
	}
	return out, nil
}
func winOverlap(index int, targets []netip.Addr) error {
	b, e := ps("@(Get-NetRoute -AddressFamily IPv4 | Select-Object DestinationPrefix,InterfaceIndex) | ConvertTo-Json -Compress")
	if e != nil {
		return e
	}
	var rows []windowsRoute
	b = []byte(strings.TrimSpace(string(b)))
	if len(b) == 0 {
		return fmt.Errorf("empty Windows route table response")
	}
	if b[0] == '{' {
		var v windowsRoute
		if e = json.Unmarshal(b, &v); e != nil {
			return e
		}
		rows = []windowsRoute{v}
	} else if e = json.Unmarshal(b, &rows); e != nil {
		return e
	}
	for _, r := range rows {
		if index != 0 && r.InterfaceIndex == index {
			continue
		}
		p, e := netip.ParsePrefix(r.DestinationPrefix)
		if e != nil {
			return e
		}
		if p.Bits() == 0 {
			continue
		}
		for _, a := range targets {
			if p.Contains(a) {
				return fmt.Errorf("overlay address %s overlaps existing route %s", a, p)
			}
		}
	}
	return nil
}
func Preflight(name string, local netip.Addr, peers []netip.Addr) error {
	if _, e := net.InterfaceByName(name); e == nil {
		return fmt.Errorf("interface %s already exists", name)
	}
	return winOverlap(0, append([]netip.Addr{local}, peers...))
}
func Configure(name string, local netip.Addr, peers []netip.Addr) error {
	var iface *net.Interface
	var e error
	for i := 0; i < 25; i++ {
		iface, e = net.InterfaceByName(name)
		if e == nil {
			break
		}
		time.Sleep(200 * time.Millisecond)
	}
	if e != nil {
		return e
	}
	idx := strconv.Itoa(iface.Index)
	if e = winOverlap(iface.Index, append([]netip.Addr{local}, peers...)); e != nil {
		return e
	}
	// Record the dedicated adapter and intended exact routes before setup,
	// so a partially successful PowerShell operation is also cleaned up.
	winOwned.Lock()
	winOwned.m[name] = ownedWin{index: iface.Index, address: local.String(), routes: append([]netip.Addr(nil), peers...)}
	winOwned.Unlock()
	var script strings.Builder
	script.WriteString("Set-NetIPInterface -InterfaceIndex " + idx + " -AddressFamily IPv4 -NlMtuBytes 1280 -Dhcp Disabled -PolicyStore ActiveStore; ")
	script.WriteString("New-NetIPAddress -InterfaceIndex " + idx + " -IPAddress '" + local.String() + "' -PrefixLength 32 -AddressFamily IPv4 -PolicyStore ActiveStore | Out-Null; ")
	for _, p := range peers {
		script.WriteString("New-NetRoute -InterfaceIndex " + idx + " -DestinationPrefix '" + p.String() + "/32' -NextHop '0.0.0.0' -RouteMetric 1 -PolicyStore ActiveStore | Out-Null; ")
	}
	if _, e = ps(script.String()); e != nil {
		return e
	}
	return nil
}
func Cleanup(name string, local netip.Addr, peers []netip.Addr) {
	winOwned.Lock()
	v, ok := winOwned.m[name]
	delete(winOwned.m, name)
	winOwned.Unlock()
	if !ok {
		return
	}
	idx := strconv.Itoa(v.index)
	var script strings.Builder
	for _, p := range v.routes {
		script.WriteString("Remove-NetRoute -InterfaceIndex " + idx + " -DestinationPrefix '" + p.String() + "/32' -PolicyStore ActiveStore -Confirm:$false -ErrorAction SilentlyContinue; ")
	}
	script.WriteString("Remove-NetIPAddress -InterfaceIndex " + idx + " -IPAddress '" + v.address + "' -PolicyStore ActiveStore -Confirm:$false -ErrorAction SilentlyContinue; ")
	_, _ = ps(script.String())
}
