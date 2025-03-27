package paillier

import (
	"math/big"

	"github.com/cronokirby/saferith"
	"github.com/taurusgroup/multi-party-sig/pkg/paillier"
	"github.com/taurusgroup/multi-party-sig/pkg/pool"
)


func GenerateKeyPair() (*paillier.PublicKey, *paillier.SecretKey, error) {
	// ダミー実装
	p := pool.NewPool(0)
	pub := &paillier.PublicKey{}
	priv := &paillier.SecretKey{}
	pub, priv = paillier.KeyGen(p)

	return pub, priv, nil
}

func Encrypt(pub *paillier.PublicKey, plaintext *big.Int) (*paillier.Ciphertext, error) {
	// ダミー実装
	m := &saferith.Int{}
	m.SetBig(plaintext, 256)
	ciphertext, nat := pub.Enc(m)
	return ciphertext, nil
}

func Decrypt(priv *paillier.SecretKey, ciphertext *paillier.Ciphertext) (*big.Int, error) {
	// ダミー実装
	m, err := priv.Dec(ciphertext)
	return m.Big(), nil
}
