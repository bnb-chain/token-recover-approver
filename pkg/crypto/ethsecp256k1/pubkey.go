package ethsecp256k1

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// NewPubKey creates a new PubKey from an ecdsa.PublicKey.
func NewPubKey(pubKey *ecdsa.PublicKey) *PubKey {
	return &PubKey{pubKey}
}

// PubKey is a wrapper around an ecdsa.PublicKey
type PubKey struct {
	pubKey *ecdsa.PublicKey
}

// Address returns the eth address of the private key.
func (pk *PubKey) Address() common.Address {
	return crypto.PubkeyToAddress(*pk.pubKey)
}

// Verify verifies a signature.
func (pk *PubKey) Verify(message []byte, signature []byte) (valid bool) {
	sig := signature[:crypto.RecoveryIDOffset] // remove recovery id
	pubKey := crypto.FromECDSAPub(pk.pubKey)
	return crypto.VerifySignature(pubKey, message, sig)
}
