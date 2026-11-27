// Static configuration entrypoint; managed products should call pkg/engine.
package main

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/xzcode/network-core-ex/internal/dataplane/relay"
	"github.com/xzcode/network-core-ex/pkg/engine"
	"golang.org/x/crypto/curve25519"
)

func run(args []string) error {
	if len(args) == 0 {
		return errors.New("usage: agent keygen | check | run | status")
	}
	switch args[0] {
	case "keygen":
		f := flag.NewFlagSet("keygen", flag.ContinueOnError)
		kind := f.String("kind", "wg", "wg or identity")
		out := f.String("out", "", "private key path (new file)")
		if e := f.Parse(args[1:]); e != nil {
			return e
		}
		if *out == "" {
			return errors.New("-out required")
		}
		if *kind == "identity" {
			pub, k, e := ed25519.GenerateKey(rand.Reader)
			if e != nil {
				return e
			}
			if e = relay.WriteSecret(*out, k); e != nil {
				return e
			}
			fmt.Printf("%x\n", pub)
			return nil
		}
		if *kind != "wg" {
			return errors.New("unknown key kind")
		}
		k := make([]byte, 32)
		if _, e := rand.Read(k); e != nil {
			return e
		}
		k[0] &= 248
		k[31] &= 127
		k[31] |= 64
		pub, e := curve25519.X25519(k, curve25519.Basepoint)
		if e != nil {
			return e
		}
		if e = relay.WriteSecret(*out, k); e != nil {
			return e
		}
		fmt.Printf("%x\n", pub)
		return nil

	case "check", "run":
		f := flag.NewFlagSet(args[0], flag.ContinueOnError)
		file := f.String("config", "", "static configuration file")
		if err := f.Parse(args[1:]); err != nil {
			return err
		}
		if *file == "" {
			return errors.New("-config required")
		}
		config, err := engine.Load(*file)
		if err != nil {
			return err
		}
		if err := engine.Check(config); err != nil {
			return err
		}
		if args[0] == "check" {
			return nil
		}
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()
		return engine.Run(ctx, config)
	case "status":
		f := flag.NewFlagSet("status", flag.ContinueOnError)
		file := f.String("file", "", "status file")
		if err := f.Parse(args[1:]); err != nil {
			return err
		}
		data, err := os.ReadFile(*file)
		if err != nil {
			return err
		}
		var status engine.Status
		if err = json.Unmarshal(data, &status); err != nil {
			return err
		}
		return json.NewEncoder(os.Stdout).Encode(status)
	default:
		return errors.New("unknown agent command")
	}
}
func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
