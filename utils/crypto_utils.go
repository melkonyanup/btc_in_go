package utils

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"log"
)

// GenerateKeyPair generates a new key pair
func GenerateKeyPair(bits int) (*rsa.PrivateKey, *rsa.PublicKey) {
	privkey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		log.Println(err)
		return nil, nil
	}
	return privkey, &privkey.PublicKey
}

// PrivateKeyToBytes private key to bytes
func PrivateKeyToBytes(priv *rsa.PrivateKey) []byte {
	privBytes := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(priv),
		},
	)

	return privBytes
}

// PublicKeyToBytes public key to bytes
func PublicKeyToBytes(pub *rsa.PublicKey) []byte {
	pubASN1, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		log.Println(err)
		return nil
	}

	pubBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubASN1,
	})

	return pubBytes
}

// BytesToPrivateKey bytes to private key
func BytesToPrivateKey(priv []byte) *rsa.PrivateKey {
	block, _ := pem.Decode(priv)
	enc := x509.IsEncryptedPEMBlock(block)
	b := block.Bytes
	var err error
	if enc {
		log.Println("is encrypted pem block")
		b, err = x509.DecryptPEMBlock(block, nil)
		if err != nil {
			log.Println(err)
			return nil
		}
	}
	key, err := x509.ParsePKCS1PrivateKey(b)
	if err != nil {
		log.Println(err)
		return nil
	}
	return key
}

// BytesToPublicKey bytes to public key
func BytesToPublicKey(pub []byte) *rsa.PublicKey {
	block, _ := pem.Decode(pub)
	enc := x509.IsEncryptedPEMBlock(block)
	b := block.Bytes
	var err error
	if enc {
		log.Println("is encrypted pem block")
		b, err = x509.DecryptPEMBlock(block, nil)
		if err != nil {
			log.Println(err)
			return nil
		}
	}
	ifc, err := x509.ParsePKIXPublicKey(b)
	if err != nil {
		log.Println(err)
		return nil
	}
	key, ok := ifc.(*rsa.PublicKey)
	if !ok {
		log.Println("Not okay when constructing public key")
		return nil
	}
	return key
}

// Hash message using SHA256
func SHA256(msg []byte) []byte {
	newhash := crypto.SHA256
	pssh := newhash.New()
	pssh.Write(msg)
	return pssh.Sum(nil)
}

// Sign a message's SHA256 digest with provided private key.
func Sign(msg []byte, sk *rsa.PrivateKey) ([]byte, error) {
	// Calculate SHA256 digest of the original message.
	newhash := crypto.SHA256
	pssh := newhash.New()
	pssh.Write(msg)
	digest := pssh.Sum(nil)

	// Sign the message.
	var opts rsa.PSSOptions
	opts.SaltLength = rsa.PSSSaltLengthAuto
	signature, err := rsa.SignPSS(rand.Reader, sk, newhash, digest, &opts)
	if err != nil {
		return nil, err
	}

	return signature, nil
}

// Verify the given signature matches the message.
func Verify(msg []byte, pk *rsa.PublicKey, signature []byte) bool {
	// Calculate SHA256 digest of the original message.
	newhash := crypto.SHA256
	pssh := newhash.New()
	pssh.Write(msg)
	digest := pssh.Sum(nil)

	var opts rsa.PSSOptions
	opts.SaltLength = rsa.PSSSaltLengthAuto

	err := rsa.VerifyPSS(pk, newhash, digest, signature, &opts)
	if err == nil {
		return true
	}
	log.Println("Fail to verify signature")
	return false
}
