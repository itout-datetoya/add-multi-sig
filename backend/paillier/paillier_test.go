package paillier

import (
	"math/big"
	"testing"
)

func TestEncryptDecrypt(t *testing.T) () {
	plaintext := big.NewInt(int64(123456789))

	pub := &PublicKey{}
	priv := &PrivateKey{}
	pub, priv, _ = GenerateKey(2048)

	ciphertext, _ := pub.Encrypt(plaintext)
	
	result, _ := priv.Decrypt(ciphertext)

	if plaintext.Int64() != result.Int64() {
		t.Errorf("got %v\nwant %v", result, plaintext)
	}
}