package ethsecp256k1

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/crypto"
)

// NewPrivKey creates a new PrivKey from a secret.
func NewPrivKey(secret string) (*PrivKey, error) {
	privKey, err := RecoveryFromPrivateKey(secret)
	if err != nil {
		return nil, err
	}
	ecdsaPrivKey, err := crypto.ToECDSA(privKey[:])
	if err != nil {
		return nil, err
	}
	return &PrivKey{ecdsaPrivKey}, nil
}

// PrivKey is a wrapper around an ecdsa.PrivateKey
type PrivKey struct {
	privKey *ecdsa.PrivateKey
}

// Sign signs a message.
func (pk *PrivKey) Sign(message []byte) (signature []byte, err error) {
	return crypto.Sign(message, pk.privKey)
}

func (pk *PrivKey) PubKey() *PubKey {
	return NewPubKey(&pk.privKey.PublicKey)
}
