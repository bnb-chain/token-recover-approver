package keymanager

import "github.com/ethereum/go-ethereum/common"

// KeyManager is the interface that wraps the basic Sign and Verify methods.
type KeyManager interface {
	Address() common.Address
	Sign(message []byte) (signature []byte, err error)
	Verify(message []byte, signature []byte) (valid bool)
}
