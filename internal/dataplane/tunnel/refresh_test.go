package tunnel

import "testing"

func TestPathRefreshPreservesMembership(t *testing.T) {
	old := []Peer{{IP: "10.160.0.2", PublicKey: "key", RelayURL: "wss://old.example/relay/v1", RelayTokenFile: "old"}}
	next := []Peer{{IP: "10.160.0.2", PublicKey: "key", RelayURL: "wss://new.example/relay/v1", RelayTokenFile: "new", ICE: true, STUNURL: "stun:new.example:3478"}}
	if !samePathMembership(old, next) {
		t.Fatal("valid transport-only refresh rejected")
	}
	next[0].PublicKey = "other"
	if samePathMembership(old, next) {
		t.Fatal("public key change accepted")
	}
	next[0].PublicKey = old[0].PublicKey
	next[0].IP = "10.160.0.3"
	if samePathMembership(old, next) {
		t.Fatal("peer IP change accepted")
	}
	next[0].IP = old[0].IP
	next[0].Subnets = []string{"192.168.0.0/24"}
	if samePathMembership(old, next) {
		t.Fatal("route change accepted")
	}
	if samePathMembership(old, nil) {
		t.Fatal("peer removal accepted")
	}
}
