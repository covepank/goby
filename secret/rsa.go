package secret

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"

	"github.com/sanbsy/gopkg/errors"
)

var (
	ErrKeyMustBePEMEncoded = errors.New("invalid key: Key must be PEM encoded")
	ErrNotRSAPrivateKey    = errors.New("key is not a valid RSA private key")
	ErrNotRSAPublicKey     = errors.New("key is not a valid RSA public key")
)

type RSAKey struct {
	privateKey *rsa.PrivateKey
	size       int
}

func CreateRSAKey(size int) (*RSAKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, size)
	if err != nil {
		return nil, err
	}
	return &RSAKey{
		privateKey: privateKey,
		size:       size,
	}, nil
}

func CreateRSAKeyByPrivateKey(privateKey *rsa.PrivateKey) *RSAKey {
	return &RSAKey{
		privateKey: privateKey,
		size:       privateKey.Size(),
	}
}
func (r *RSAKey) PrivateKey() *rsa.PrivateKey {
	return r.privateKey
}

func (r *RSAKey) PublicKey() *rsa.PublicKey {
	return &r.privateKey.PublicKey
}

func (r *RSAKey) Size() int {
	return r.size
}
func (r *RSAKey) EncodePrivateToPEM() ([]byte, error) {
	pkcs, err := x509.MarshalPKCS8PrivateKey(r.privateKey)
	if err != nil {
		return nil, err
	}
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: pkcs,
	}
	return pem.EncodeToMemory(block), nil
}
func (r *RSAKey) EncodePublicToPEM() ([]byte, error) {
	pub, err := x509.MarshalPKIXPublicKey(r.PublicKey())
	if err != nil {
		return nil, err
	}
	block := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pub,
	}
	return pem.EncodeToMemory(block), nil
}
func ParseRSAPublicKey(data []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, ErrKeyMustBePEMEncoded
	}

	parsedKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		if cert, err := x509.ParseCertificate(block.Bytes); err == nil {
			parsedKey = cert.PublicKey
		} else {
			return nil, err
		}
	}

	if pkey, ok := parsedKey.(*rsa.PublicKey); !ok {
		return nil, ErrNotRSAPublicKey
	} else {
		return pkey, nil
	}

}

func ParseRSAPrivateKey(key []byte) (*rsa.PrivateKey, error) {
	var err error

	var block *pem.Block
	if block, _ = pem.Decode(key); block == nil {
		return nil, ErrKeyMustBePEMEncoded
	}

	var parsedKey interface{}
	if parsedKey, err = x509.ParsePKCS1PrivateKey(block.Bytes); err != nil {
		if parsedKey, err = x509.ParsePKCS8PrivateKey(block.Bytes); err != nil {
			return nil, err
		}
	}

	pkey, ok := parsedKey.(*rsa.PrivateKey)
	if !ok {
		return nil, ErrNotRSAPrivateKey
	}

	return pkey, nil
}
