package secret

import (
	"crypto/rsa"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
)

func initRSAKEY() (*rsa.PublicKey, *rsa.PrivateKey, error) {
	pubCont, err := ioutil.ReadFile("./rsa_public_key.pem")
	if err != nil {
		return nil, nil, err
	}
	pub, err := ParseRSAPublicKey(pubCont)
	if err != nil {
		return nil, nil, err
	}
	priCont, err := ioutil.ReadFile("./rsa_private_key.pem")
	if err != nil {
		return nil, nil, err
	}
	pri, err := ParseRSAPrivateKey(priCont)
	if err != nil {
		return nil, nil, err
	}
	return pub, pri, err
}

func TestSigner_Sign(t *testing.T) {
	data := []byte("hello world123132132156132asdasdadsasdfasdfsd")
	pub, pri, err := initRSAKEY()
	if err != nil {
		t.Error(err)
		return
	}
	signature, err := RSASign(data, SigningMethodRSA512, pri)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(signature)
	err = RSAVerify(data, signature, SigningMethodRSA512, pub)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestAesEncrypt(t *testing.T) {
	data := []byte(`《山海经》是中国一部记述古代志怪的古籍，大体是战国中后期到汉代初中期的楚国或巴蜀人所作。也是一部荒诞不经的奇书。
该书作者不详，古人认为该书是“战国好奇之士取《穆王传》，杂录《庄》、《列》 、《离骚》 、《周书》、《晋乘》以成者” 。
现代学者也均认为成书并非一时，作者亦非一人。`)
	key := []byte("hello world! bei")
	crypted, err := AESEncrypt(data, key)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(crypted)

	decrypted, err := AESDecrypt(crypted, key)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(string(decrypted))

}
func TestRsaEncrypt(t *testing.T) {
	data := []byte(`1024102412123151215`)
	pub, pri, err := initRSAKEY()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(pub.Size())
	encrypted, err := RSAEncrypt(data, pub)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(encrypted)
	origin, err := RSADecrypt(encrypted, pri)
	if err != nil {
		t.Log(err)
		return
	}
	t.Log(string(origin))
}

func TestECDSASign(t *testing.T) {
	r := require.New(t)

	content := []byte(`hello world!!!`)
	ecdsaKey, err := CreateECDSAKey()
	r.Nil(err)

	signature, err := ECDSASign(content, SigningMethodRSA512, ecdsaKey.PrivateKey())
	r.Nil(err)
	t.Log(signature)

	err = ECDSAVerify(content, signature, SigningMethodRSA512, ecdsaKey.PublicKey())
	r.Nil(err)
}
