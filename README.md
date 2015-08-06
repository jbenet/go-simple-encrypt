# go-simple-encrypt

Trivial tool to {en,de}crypt things symmetrically. This is not meant to be very strong, just better than storing things in the clear. It uses AES in CTR mode, generates a random IV, and stores it at the beginning of the ciphertext.

**Security Warning:** please don't use this for anything serious. It should be correct, but it's not audited.

TODO: maybe use AEAD with `cipher.NewGCM`

## Usage of senc tool

### Install

Install with `go get`:

```sh
go get -u github.com/jbenet/go-simple-encrypt/senc
```

Or get a prebuilt binary from gobuilder:

> https://gobuilder.me/github.com/jbenet/go-simple-encrypt/senc

### Generate a key

AES256 keys are just 256 bits of randomness. We use them encoded in base58.

```sh
> key=$(senc --key-gen)
> echo $key
GjEcPVFWZCoU31LBJSrNnLwqzv4biZ3BJTrP9ddGQ63

# this is the same as:
> key=$(head -c 32 /dev/urandom | base58)
> echo $key
3uoGkxdJrSmWvw3MnbtNhNRE6hCDp7kBd9hXBBG5Doiw
```

### Encrypt

```sh
cat plaintext | senc -k $key -e >ciphertext
```

### Decrypt

```sh
cat ciphertext | senc -k $key -d >plaintext
```

### Test it

```sh
key=$(senc --key-gen)
head -c 1048576 /dev/urandom >plaintext
senc -k $key -e <plaintext >ciphertext
senc -k $key -d <ciphertext >plaintext2
diff plaintext plaintext2
```

## Usage of lib

#### Godoc: https://godoc.org/github.com/jbenet/go-simple-encrypt

```
import (
  senc "github.com/jbenet/go-simple-encrypt"
)

func encrypt() {
  k := getKey() // 256bit key. can use: senc.RandomKey()

  plaintext := io.Stdin
  ciphertext, err := senc.Encrypt(plaintext, k)
  if err != nil {
    panic(err)
  }

  _, err = io.Copy(io.Stdout, ciphertext)
  if err != nil { // dont forget to check this!
    panic(err)
  }
}

func decrypt() {
  k := getKey() // 256bit key

  ciphertext := io.Stdin
  plaintext, err := senc.Decrypt(ciphertext, k)
  if err != nil {
    panic(err)
  }

  _, err = io.Copy(io.Stdout, ciphertext)
  if err != nil { // dont forget to check this!
    panic(err)
  }
}
```
