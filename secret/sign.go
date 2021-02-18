package secret

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"math/big"

	"github.com/sanbsy/goby/errors"
)

const (
	ErrBase64DecodeError    = "签名base64解码异常"
	ErrSignatureCheckError  = "签名校验失败"
	ErrSignatureCreateError = "签名创建失败"
	ErrSigningMethodError   = "错误的签名方法"
)

const (
	SigningMethodRSA256 = "SHA256"
	SigningMethodRSA512 = "SHA512"
	SigningMethodRSA384 = "SHA384"
)

// SHA Hash
var signingMethodMap = map[string]crypto.Hash{
	SigningMethodRSA256: crypto.SHA256,
	SigningMethodRSA384: crypto.SHA384,
	SigningMethodRSA512: crypto.SHA512,
}

// RSASign 使用 RSA 签名
func RSASign(data []byte, method string, privateKey *rsa.PrivateKey) (string, error) {
	hashed, err := hashSum(data, method)
	if err != nil {
		return "", errors.Wrap(err, ErrSigningMethodError)
	}

	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, signingMethodMap[method], hashed)
	if err != nil {
		return "", errors.Wrap(err, ErrSignatureCreateError)
	}
	return base64.StdEncoding.EncodeToString(signature), nil
}

// RSAVerify rsa验证签名
func RSAVerify(content []byte, signature string, method string, publicKey *rsa.PublicKey) error {
	hashed, err := hashSum(content, method)
	if err != nil {
		return errors.Wrap(err, ErrSigningMethodError)
	}
	// base64 解码 签名
	sign, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return errors.Wrap(err, ErrBase64DecodeError)
	}

	// 返回验签结果
	err = rsa.VerifyPKCS1v15(publicKey, signingMethodMap[method], hashed[:], sign)
	if err != nil {
		return errors.Wrap(err, ErrSignatureCheckError)
	}
	return nil
}

// ECDSASign ecdsa签名
func ECDSASign(data []byte, method string, privateKey *ecdsa.PrivateKey) (string, error) {
	hashed, err := hashSum(data, method)
	if err != nil {
		return "", errors.Wrap(err, ErrSigningMethodError)
	}

	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hashed)
	rt, err := r.MarshalText()
	if err != nil {
		return "", errors.Wrap(err, ErrSignatureCreateError)
	}
	st, err := s.MarshalText()
	if err != nil {
		return "", errors.Wrap(err, ErrSignatureCreateError)
	}

	buf := new(bytes.Buffer)
	buf.Write(rt)
	buf.WriteByte('.')
	buf.Write(st)
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil

}

// ECDSAVerify ecdsa 验证签名
func ECDSAVerify(content []byte, signature string, method string, publicKey *ecdsa.PublicKey) error {
	hashed, err := hashSum(content, method)
	if err != nil {
		return errors.Wrap(err, ErrSigningMethodError)
	}
	signatureBytes, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return errors.Wrap(err, ErrBase64DecodeError)
	}
	r, s, err := parseECDSASignature(signatureBytes)
	if err != nil {
		return errors.Wrap(err, ErrSignatureCheckError)
	}

	cr := ecdsa.Verify(publicKey, hashed, r, s)
	if !cr {
		return errors.New(ErrSignatureCheckError)
	}
	return nil

}

// HMACSHA256  hmac 签名
func HMACSHA256(data, key []byte) string {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// HMACSHA512 hmac 验签
func HMACSHA512(data, key []byte) string {
	h := hmac.New(sha512.New, key)
	h.Write(data)
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func hashSum(data []byte, method string) ([]byte, error) {
	var hashed []byte
	if _, ok := signingMethodMap[method]; !ok {
		return nil, errors.New("invalid singing method.")
	}
	switch method {
	case SigningMethodRSA256:
		h := sha256.Sum256(data)
		hashed = h[:]
	case SigningMethodRSA384:
		h := sha512.Sum384(data)
		hashed = h[:]
	case SigningMethodRSA512:
		h := sha512.Sum512(data)
		hashed = h[:]
	default:
		return nil, errors.New("invalid signing method.")
	}
	return hashed, nil
}

func parseECDSASignature(data []byte) (*big.Int, *big.Int, error) {
	r, s := new(big.Int), new(big.Int)

	idx := 0
	for i, v := range data {
		if v == '.' {
			idx = i
			break
		}
	}
	if idx <= 0 || idx >= (len(data)-1) {
		return nil, nil, errors.New("invalid signature")
	}

	if err := r.UnmarshalText(data[:idx]); err != nil {
		return nil, nil, nil
	}
	if err := s.UnmarshalText(data[idx+1:]); err != nil {
		return nil, nil, err
	}
	return r, s, nil
}
