package transport

import (
	"net/netip"
	"testing"
)

func TestReplaceClosed(t *testing.T) {
	peer := Peer{Logical: netip.MustParseAddrPort("10.160.0.2:1"), Direct: netip.MustParseAddrPort("127.0.0.1:10001")}
	b, e := New([]Peer{peer})
	if e != nil {
		t.Fatal(e)
	}
	if _, _, e = b.Open(0); e != nil {
		t.Fatal(e)
	}
	defer b.Close()
	next := peer
	next.Direct = netip.MustParseAddrPort("127.0.0.1:10002")
	if e = b.ReplaceClosed([]Peer{next}); e == nil {
		t.Fatal("live bind replacement accepted")
	}
	if e = b.Close(); e != nil {
		t.Fatal(e)
	}
	if e = b.ReplaceClosed([]Peer{next}); e != nil {
		t.Fatal(e)
	}
	if _, _, e = b.Open(0); e != nil {
		t.Fatal(e)
	}
	if _, e = b.ParseEndpoint(next.Logical.String()); e != nil {
		t.Fatal(e)
	}
}
