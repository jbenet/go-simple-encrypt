package senc

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

const KeyLength = 32 // aes256

func readIV(r io.Reader) ([]byte, error) {
	iv := make([]byte, aes.BlockSize)
	_, err := io.ReadFull(r, iv)
	return iv, err
}

func RandomKey() ([]byte, error) {
	k := make([]byte, KeyLength)
	_, err := io.ReadFull(rand.Reader, k)
	return k, err
}

func Encrypt(key []byte, plaintext io.Reader) (ciphertext io.Reader, err error) {
	bc, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	iv, err := readIV(rand.Reader) // generate random IV
	if err != nil {
		return nil, err
	}
	ivr := bytes.NewReader(iv)

	stream := cipher.NewCTR(bc, iv)
	sr := &cipher.StreamReader{S: stream, R: plaintext}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	return io.MultiReader(ivr, sr), nil
}

func Decrypt(key []byte, ciphertext io.Reader) (plaintext io.Reader, err error) {
	bc, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	iv, err := readIV(ciphertext)
	if err != nil {
		return nil, err
	}

	stream := cipher.NewCTR(bc, iv)
	sr := &cipher.StreamReader{S: stream, R: ciphertext}
	return sr, nil
}
