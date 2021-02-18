package secret

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

type ECDSAKey struct {
	private *ecdsa.PrivateKey
}

var (
	// ErrNotECDSAPrivateKey = errors.New("key is not a valid ECDSA private key")
	ErrNotECDSAPublicKey = errors.New("key is not a valid ECDSA public key")
)

func CreateECDSAKey() (*ECDSAKey, error) {
	curve := elliptic.P521()
	pri, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return nil, err
	}
	return &ECDSAKey{
		private: pri,
	}, nil
}

func CreateECDSAKeyByPrivateKey(key *ecdsa.PrivateKey) *ECDSAKey {
	return &ECDSAKey{private: key}
}

func (ek *ECDSAKey) EncodePrivateToPEM() ([]byte, error) {
	pri, err := x509.MarshalECPrivateKey(ek.private)
	if err != nil {
		return nil, err
	}
	block := &pem.Block{
		Type:  "ECDSA PRIVATE KEY",
		Bytes: pri,
	}
	return pem.EncodeToMemory(block), nil
}

func (ek *ECDSAKey) EncodePublicToPem() ([]byte, error) {
	pub, err := x509.MarshalPKIXPublicKey(&(ek.private.PublicKey))
	if err != nil {
		return nil, err
	}

	block := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pub,
	}
	return pem.EncodeToMemory(block), nil
}

func (ek *ECDSAKey) PublicKey() *ecdsa.PublicKey {
	return &ek.private.PublicKey
}

func (ek *ECDSAKey) PrivateKey() *ecdsa.PrivateKey {
	return ek.private
}

func ParseECDSAPublic(data []byte) (*ecdsa.PublicKey, error) {
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, ErrKeyMustBePEMEncoded
	}

	publicStream, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	public, ok := publicStream.(*ecdsa.PublicKey)
	if !ok {
		return nil, ErrNotECDSAPublicKey
	}
	return public, nil
}

func ParseECDSAPrivate(data []byte) (*ecdsa.PrivateKey, error) {
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, ErrKeyMustBePEMEncoded
	}

	return x509.ParseECPrivateKey(block.Bytes)
}
