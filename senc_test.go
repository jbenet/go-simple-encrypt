package senc

import (
	"bytes"
	"crypto/aes"
	"crypto/rand"
	"io/ioutil"
	"testing"
)

func randKey(t *testing.T) []byte {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		t.Fatal(err)
	}
	return key
}

func TestEncrypt(t *testing.T) {
	key := randKey(t)
	ptx := make([]byte, 1048576)

	if _, err := rand.Read(ptx); err != nil {
		t.Fatal(err)
	}

	ctxr, err := Encrypt(key, bytes.NewReader(ptx))
	if err != nil {
		t.Fatal("error during encryption", err)
	}

	ctx, err := ioutil.ReadAll(ctxr)
	if err != nil {
		t.Fatal("reading error", err)
	}

	if len(ctx) != len(ptx)+aes.BlockSize {
		t.Fatal("ciphertext of wrong length")
	}

	if bytes.Equal(ptx, ctx[aes.BlockSize:]) {
		t.Fatal("no encryption")
	}

	ptx2r, err := Decrypt(key, bytes.NewReader(ctx))
	if err != nil {
		t.Fatal("error during decryption", err)
	}

	ptx2, err := ioutil.ReadAll(ptx2r)
	if err != nil {
		t.Fatal("reading error", err)
	}

	if !bytes.Equal(ptx, ptx2) {
		t.Error("plaintexts differ")
	}
}
