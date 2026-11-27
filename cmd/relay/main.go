// Standalone authenticated relay/signaling service. No management API.
package main

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/tls"
	"errors"
	"flag"
	"fmt"
	"github.com/xzcode/network-core-ex/internal/dataplane/relay"
	"github.com/xzcode/network-core-ex/pkg/stun"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func run(args []string) error {
	if len(args) > 0 {
		switch args[0] {
		case "issuer-keygen":
			f := flag.NewFlagSet("issuer-keygen", flag.ContinueOnError)
			out := f.String("out", "", "private key path")
			if e := f.Parse(args[1:]); e != nil {
				return e
			}
			if *out == "" {
				return errors.New("-out required")
			}
			pub, k, e := ed25519.GenerateKey(rand.Reader)
			if e != nil {
				return e
			}
			if e = relay.WriteSecret(*out, k); e != nil {
				return e
			}
			fmt.Printf("%x\n", pub)
			return nil
		case "lease":
			f := flag.NewFlagSet("lease", flag.ContinueOnError)
			key := f.String("issuer-key", "", "private key file")
			node := f.String("node", "", "node ID")
			network := f.String("network", "", "network ID")
			flow := f.String("flow", "", "flow ID")
			a := f.String("a", "", "device A identity public hex")
			b := f.String("b", "", "device B identity public hex")
			ttl := f.Duration("ttl", 10*time.Minute, "max 15m")
			out := f.String("out", "", "new token file")
			if e := f.Parse(args[1:]); e != nil {
				return e
			}
			if *out == "" || *ttl <= 0 || *ttl > 15*time.Minute {
				return errors.New("-out and valid ttl required")
			}
			k, e := relay.ReadHex(*key, ed25519.PrivateKeySize, true)
			if e != nil {
				return e
			}
			v, e := relay.Issue(k, relay.Lease{Node: *node, Network: *network, Flow: *flow, A: *a, B: *b, Expires: time.Now().Add(*ttl).Unix()})
			if e != nil {
				return e
			}
			file, e := os.OpenFile(*out, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
			if e != nil {
				return e
			}
			defer file.Close()
			_, e = file.WriteString(v + "\n")
			return e
		}
	}
	f := flag.NewFlagSet("relay", flag.ContinueOnError)
	listen := f.String("listen", "127.0.0.1:8443", "TLS listener")
	node := f.String("node", "", "node identity")
	issuer := f.String("issuer-public", "", "issuer public key file")
	cert := f.String("tls-cert", "", "TLS certificate")
	key := f.String("tls-key", "", "TLS key")
	stunListen := f.String("stun-listen", "", "optional STUN UDP4 listener")
	if err := f.Parse(args); err != nil {
		return err
	}
	if *node == "" || *cert == "" || *key == "" {
		return errors.New("node, tls-cert and tls-key required")
	}
	pub, err := relay.ReadHex(*issuer, ed25519.PublicKeySize, false)
	if err != nil {
		return err
	}
	pair, err := tls.LoadX509KeyPair(*cert, *key)
	if err != nil {
		return err
	}
	listener, err := net.Listen("tcp", *listen)
	if err != nil {
		return err
	}
	defer listener.Close()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	mux := http.NewServeMux()
	mux.Handle("/relay/v1", relay.NewServer(*node, pub))
	mux.Handle("/signal/v1", relay.NewSignalServer(*node, pub))
	srv := &http.Server{Handler: mux, ReadHeaderTimeout: 5 * time.Second, IdleTimeout: 30 * time.Second, TLSConfig: &tls.Config{MinVersion: tls.VersionTLS12, Certificates: []tls.Certificate{pair}}}
	errorsCh := make(chan error, 2)
	if *stunListen != "" {
		addr, err := net.ResolveUDPAddr("udp4", *stunListen)
		if err != nil {
			return err
		}
		c, err := net.ListenUDP("udp4", addr)
		if err != nil {
			return err
		}
		defer c.Close()
		go func() { errorsCh <- stun.Serve(ctx, c) }()
	}
	go func() { errorsCh <- srv.ServeTLS(listener, "", "") }()
	select {
	case <-ctx.Done():
		shutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return srv.Shutdown(shutdown)
	case err := <-errorsCh:
		if errors.Is(err, http.ErrServerClosed) || ctx.Err() != nil {
			return nil
		}
		return err
	}
}
func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
