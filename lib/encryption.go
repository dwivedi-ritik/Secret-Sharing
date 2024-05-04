package lib

//Shamelessy copied from this gist https://gist.github.com/miguelmota/3ea9286bd1d3c2a985b67cac4ba2130a
//Credit: https://gist.github.com/miguelmota

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/pem"
	"log"
)

type EncryptionKey struct {
	PrivateKey []byte
	PublicKey  []byte
}

type PublicKey struct {
	Key []byte
}

type RSAEncryption struct {
	Keys EncryptionKey
}

func (rsaEncryption *RSAEncryption) EncryptMessage(msg []byte) []byte {
	hash := sha512.New()
	ciphertext, err := rsa.EncryptOAEP(hash, rand.Reader, rsaEncryption.bytesToPublicKey(), msg, nil)
	if err != nil {
		panic(err)
	}
	return ciphertext
}

func (rsaEncryption *RSAEncryption) DecryptMessage(ciphertext []byte) []byte {
	hash := sha512.New()
	plaintext, err := rsa.DecryptOAEP(hash, rand.Reader, rsaEncryption.bytesToPrivateKey(), ciphertext, nil)
	if err != nil {
		panic(err)
	}
	return plaintext
}

func (rsaEncryption *RSAEncryption) bytesToPublicKey() *rsa.PublicKey {
	block, _ := pem.Decode(rsaEncryption.Keys.PublicKey)
	ifc, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		panic(err)
	}
	key, ok := ifc.(*rsa.PublicKey)
	if !ok {
		log.Fatal("Public Key Invalid")
	}
	return key
}

func (rsaEncryption *RSAEncryption) bytesToPrivateKey() *rsa.PrivateKey {
	block, _ := pem.Decode(rsaEncryption.Keys.PrivateKey)
	ifc, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		panic(err)
	}

	return ifc
}

func generateRSAKeyPair(bits int) *EncryptionKey {
	privkey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		panic(err)
	}
	privBytes := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privkey),
		},
	)
	pubASN1, err := x509.MarshalPKIXPublicKey(&privkey.PublicKey)
	if err != nil {
		panic(err)
	}

	pubBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubASN1,
	})
	return &EncryptionKey{PrivateKey: privBytes, PublicKey: pubBytes}
}

func NewRSAEncryption() *RSAEncryption {
	const bitsSize = 2048
	encryptionKey := generateRSAKeyPair(bitsSize)

	return &RSAEncryption{Keys: *encryptionKey}

}
