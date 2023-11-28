package local

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/bnb-chain/token-recover-approver/pkg/crypto/ethsecp256k1"
	"github.com/bnb-chain/token-recover-approver/pkg/keymanager"
)

var _ keymanager.KeyManager = (*LocalKeyManager)(nil)

func NewLocalKeyManager(secret string) (*LocalKeyManager, error) {
	privKey, err := ethsecp256k1.NewPrivKey(secret)
	if err != nil {
		return nil, err
	}
	return &LocalKeyManager{privKey}, nil
}

type LocalKeyManager struct {
	privKey *ethsecp256k1.PrivKey
}

// Address implements keymanager.KeyManager.
func (km *LocalKeyManager) Address() common.Address {
	return km.privKey.PubKey().Address()
}

// Sign implements keymanager.KeyManager.
func (km *LocalKeyManager) Sign(message []byte) (signature []byte, err error) {
	return km.privKey.Sign(message)
}

// Verify implements keymanager.KeyManager.
func (km *LocalKeyManager) Verify(message []byte, signature []byte) (valid bool) {
	return km.privKey.PubKey().Verify(message, signature)
}
