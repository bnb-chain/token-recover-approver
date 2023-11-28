package aws

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/bnb-chain/token-recover-approver/pkg/crypto/ethsecp256k1"
	"github.com/bnb-chain/token-recover-approver/pkg/keymanager"
)

var _ keymanager.KeyManager = (*SecretManager)(nil)

func NewSecretManager(secretName, region string) (*SecretManager, error) {
	secretString, err := getAWSSecretString(secretName, region)
	if err != nil {
		return nil, err
	}
	privKey, err := ethsecp256k1.NewPrivKey(secretString)
	if err != nil {
		return nil, err
	}
	return &SecretManager{privKey}, nil
}

type SecretManager struct {
	privKey *ethsecp256k1.PrivKey
}

// Address implements keymanager.KeyManager.
func (km *SecretManager) Address() common.Address {
	return km.privKey.PubKey().Address()
}

// Sign implements keymanager.KeyManager.
func (km *SecretManager) Sign(message []byte) (signature []byte, err error) {
	return km.privKey.Sign(message)
}

// Verify implements keymanager.KeyManager.
func (km *SecretManager) Verify(message []byte, signature []byte) (valid bool) {
	return km.privKey.PubKey().Verify(message, signature)
}
