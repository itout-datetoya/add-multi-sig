package paillier

import (
	"math/big"
	"encoding/json"
	"crypto/rand"
	"errors"
    pb "multisigservice/proto/paillierpb"
)

// L(u) = (u - 1) / n
func L(u, n *big.Int) *big.Int {
    res := new(big.Int).Sub(u, big.NewInt(1))
    return res.Div(res, n)
}

type PublicKey struct {
    N       *big.Int
    NSquare *big.Int
    G       *big.Int
}

// PublicKey → Protobuf
func (pk *PublicKey) ToProto() *pb.PublicKey {
    return &pb.PublicKey{
        N:       pk.N.Bytes(),
        NSquare: pk.NSquare.Bytes(),
        G:       pk.G.Bytes(),
    }
}

// Protobuf → PublicKey
func PublicKeyFromProto(msg *pb.PublicKey) *PublicKey {
    return &PublicKey{
        N:       new(big.Int).SetBytes(msg.N),
        NSquare: new(big.Int).SetBytes(msg.NSquare),
        G:       new(big.Int).SetBytes(msg.G),
    }
}

type PrivateKey struct {
    PublicKey
    Lambda   *big.Int
    Mu       *big.Int
}

// PrivateKey → Protobuf
func (sk *PrivateKey) ToProto() *pb.PrivateKey {
    return &pb.PrivateKey{
        PublicKey: sk.PublicKey.ToProto(),
        Lambda:    sk.Lambda.Bytes(),
        Mu:        sk.Mu.Bytes(),
    }
}

// Protobuf → PrivateKey
func PrivateKeyFromProto(msg *pb.PrivateKey) *PrivateKey {
    return &PrivateKey{
        PublicKey: *PublicKeyFromProto(msg.PublicKey),
        Lambda:    new(big.Int).SetBytes(msg.Lambda),
        Mu:        new(big.Int).SetBytes(msg.Mu),
    }
}

type Ciphertext struct {
	c *big.Int
}


// Ciphertext → Protobuf
func (ct *Ciphertext) ToProto() *pb.Ciphertext {
    return &pb.Ciphertext{
        C: ct.c.Bytes(),
    }
}

// Protobuf → Ciphertext
func CiphertextFromProto(msg *pb.Ciphertext) *Ciphertext {
    return &Ciphertext{
        c: new(big.Int).SetBytes(msg.C),
    }
}

// GenerateKey generates a Paillier keypair of specified bit length (e.g., 2048 or 3072)
func GenerateKey(bitLen int) (*PublicKey, *PrivateKey, error) {
    // 1. 素数p, qの生成
    p, err := rand.Prime(rand.Reader, bitLen/2)
    if err != nil {
        return nil, nil, err
    }
    q, err := rand.Prime(rand.Reader, bitLen/2)
    if err != nil {
        return nil, nil, err
    }
    n := new(big.Int).Mul(p, q)
    nSquare := new(big.Int).Mul(n, n)

    // 2. λ = lcm(p-1, q-1)
    p1 := new(big.Int).Sub(p, big.NewInt(1))
    q1 := new(big.Int).Sub(q, big.NewInt(1))
    gcd := new(big.Int).GCD(nil, nil, p1, q1)
    lambda := new(big.Int).Div(new(big.Int).Mul(p1, q1), gcd)

    // 3. g = n + 1
    g := new(big.Int).Add(n, big.NewInt(1))

    // 4. μ = (L(g^λ mod n^2))^-1 mod n
    //    m = g^λ mod n^2
    m := new(big.Int).Exp(g, lambda, nSquare)
    l := L(m, n)
    mu := new(big.Int).ModInverse(l, n)
    if mu == nil {
        return nil, nil, errors.New("failed to compute modular inverse for mu")
    }

    // 公開鍵部分を構築
    pub := &PublicKey{
        N:       n,
        NSquare: nSquare,
        G:       g,
    }

    // 秘密鍵部分を構築（PublicKey を埋め込み）
    priv := &PrivateKey{
        PublicKey: *pub,
        Lambda:    lambda,
        Mu:        mu,
    }

    return pub, priv, nil
}

// Encrypt encrypts plaintext m ∈ [0, n)
func (pub *PublicKey) Encrypt(m *big.Int) (*Ciphertext, error) {
    if m.Cmp(pub.N) >= 0 {
        return nil, errors.New("plaintext too large")
    }
    // 1. 乱数 r ∈ [1, n), gcd(r,n)=1
    r, err := rand.Prime(rand.Reader, pub.N.BitLen())
    if err != nil {
        return nil, err
    }
    // 2. c = g^m * r^n mod n^2
    gm := new(big.Int).Exp(pub.G, m, pub.NSquare)
    rn := new(big.Int).Exp(r, pub.N, pub.NSquare)
    ct := &Ciphertext{new(big.Int).Mod(new(big.Int).Mul(gm, rn), pub.NSquare)}
    return ct, nil
}

// Decrypt recovers plaintext from ciphertext c
func (priv *PrivateKey) Decrypt(ct *Ciphertext) (*big.Int, error) {
    if ct.c.Cmp(priv.NSquare) >= 0 {
        return nil, errors.New("ciphertext too large")
    }
    // m = L(c^λ mod n^2) * μ mod n
    u := new(big.Int).Exp(ct.c, priv.Lambda, priv.NSquare)
    l := L(u, priv.N)
    m := new(big.Int).Mod(new(big.Int).Mul(l, priv.Mu), priv.N)
    return m, nil
}

func (ct *Ciphertext) AddScalar(pub *PublicKey, m *big.Int) (*Ciphertext, error) {
    if m.Sign() < 0 {
        return nil, errors.New("scalar must be non-negative")
    }
    // g^m mod n^2
    gm := new(big.Int).Exp(pub.G, m, pub.NSquare)
    // c' = c * g^m mod n^2
    cNew := new(big.Int).Mod(new(big.Int).Mul(ct.c, gm), pub.NSquare,)
    return &Ciphertext{c: cNew}, nil
}

func (ct *Ciphertext) Add(pub *PublicKey, ct2 *Ciphertext) (*Ciphertext, error) {
    // c' = c1 * c2 mod n^2
    cNew := new(big.Int).Mod(new(big.Int).Mul(ct.c, ct2.c), pub.NSquare)
    return &Ciphertext{c: cNew}, nil
}

func (ct *Ciphertext) MulScalar(pub *PublicKey, k *big.Int) (*Ciphertext, error) {
    if k.Sign() < 0 {
        return nil, errors.New("scalar must be non-negative")
    }
    // c' = c^k mod n^2
    cNew := new(big.Int).Exp(ct.c, k, pub.NSquare)
    return &Ciphertext{c: cNew}, nil
}