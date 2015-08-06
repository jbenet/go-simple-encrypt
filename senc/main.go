package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	b58 "github.com/jbenet/go-base58"

	senc "github.com/jbenet/go-simple-encrypt"
)

type options struct {
	encrypt bool
	decrypt bool
	keygen  bool
	keyS    string
	key     []byte
}

func onlyOne(cs ...bool) bool {
	var c int
	for _, a := range cs {
		if a {
			c++
		}
	}
	return c == 1
}

func parseOpts() (options, error) {
	o := options{}

	flag.StringVar(&o.keyS, "k", "", "key to use (in base58)")
	flag.BoolVar(&o.encrypt, "e", false, "encrypt mode")
	flag.BoolVar(&o.decrypt, "d", false, "decrypt mode")
	flag.BoolVar(&o.keygen, "key-gen", false, "generate a key")

	flag.Usage = usage
	flag.Parse()

	if !onlyOne(o.keygen, o.encrypt, o.decrypt) {
		return o, fmt.Errorf("must choose either -e or -d or --key-gen")
	}

	if !o.keygen {
		o.keyS = strings.TrimSpace(o.keyS)
		o.key = b58.Decode(o.keyS)
		if o.keyS == "" || o.key == nil {
			return o, fmt.Errorf("must provide a key in base58 with -k")
		}

		if len(o.key) < senc.KeyLength {
			return o, fmt.Errorf("key too short. must be 256 bits, decoded.")
		} else if len(o.key) > senc.KeyLength {
			return o, fmt.Errorf("key too short. must be 256 bits, decoded.")
		}
	}

	return o, nil
}

func usage() {
	p := os.Args[0]
	fmt.Println("usage: ", p, " -e -k <key-in-base58> - encrypt stdin with aes")
	fmt.Println("       ", p, " -d -k <key-in-base58> - decrypt stdin with aes")
	fmt.Println("")
	fmt.Println("options")
	flag.PrintDefaults()
	os.Exit(0)
}

func run() error {
	opts, err := parseOpts()
	if err != nil {
		return err
	}

	if opts.keygen {
		k, err := senc.RandomKey()
		if err != nil {
			return err
		}
		fmt.Println(b58.Encode(k))
		return nil
	}

	var r io.Reader
	switch {
	case opts.encrypt:
		r, err = senc.Encrypt(opts.key, os.Stdin)
	case opts.decrypt:
		r, err = senc.Decrypt(opts.key, os.Stdin)
	default:
		return fmt.Errorf("must choose either -e or -d")
	}
	if err != nil {
		return err
	}

	_, err = io.Copy(os.Stdout, r)
	return err
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
