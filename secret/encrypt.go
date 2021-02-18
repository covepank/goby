package secret

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
)

// AESEncrypt AES CBC加密
// key 长度必须为 16， 24 或 32
func AESEncrypt(origin []byte, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	blockSize := block.BlockSize()
	// 补全码
	data := PKCS7Padding(origin, blockSize)
	// 加密模式
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	encrypted := make([]byte, len(data))
	blockMode.CryptBlocks(encrypted, data)
	return base64.StdEncoding.EncodeToString(encrypted), nil
}
func AESDecrypt(encrypted string, key []byte) ([]byte, error) {
	encryptedByte, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	// 加密模式
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	originData := make([]byte, len(encryptedByte))
	blockMode.CryptBlocks(originData, encryptedByte)
	// 去补全码
	return PKCS7UnPadding(originData), nil
}

//  RSA加密
// 待加密字段长度不能大于 117 （128-11）
func RSAEncrypt(origData []byte, publicKey *rsa.PublicKey) (string, error) {
	encrypted, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, origData)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(encrypted), nil
}

// 解密
func RSADecrypt(encrypted string, privateKey *rsa.PrivateKey) ([]byte, error) {
	encryptedByte, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return nil, err
	}
	return rsa.DecryptPKCS1v15(rand.Reader, privateKey, encryptedByte)
}

// 补码
func PKCS7Padding(text []byte, size int) []byte {
	padding := size - len(text)%size
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(text, padtext...)
}

// 去补码
func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
