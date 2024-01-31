package memory

import (
	"github.com/cosmos/cosmos-sdk/types"

	"github.com/bnb-chain/token-recover-approver/internal/store"
	"github.com/bnb-chain/token-recover-approver/pkg/util"
)

var _ store.Store = (*MemoryStore)(nil)

func NewMemoryStore(proofsPath string) (*MemoryStore, error) {
	errChan := make(chan error, 1)
	defer close(errChan)
	// load proofs
	stream := util.NewJSONStream(func() any {
		return &Proof{}
	})

	proofs := make(map[string]*Proof)
	go func() {
		for data := range stream.Watch() {
			if data.Error != nil {
				errChan <- data.Error
				return
			}
			proof := data.Data.(*Proof)
			index := proof.Address.String() + ":" + proof.Coin.Denom
			proofs[index] = proof
		}
		errChan <- nil
	}()
	stream.Start(proofsPath)
	err := <-errChan
	if err != nil {
		return nil, err
	}

	return &MemoryStore{
		proofs: proofs,
	}, nil
}

// MemoryStore implements store.Store.
type MemoryStore struct {
	proofs map[string]*Proof // address:index:symbol -> proofs
}

// GetAccountProofs implements store.Store.
func (ss *MemoryStore) GetAccountAssetProof(address types.AccAddress, symbol string) (*store.Proof, error) {

	index := address.String() + ":" + symbol
	proofs, exist := ss.proofs[index]
	if !exist {
		return nil, ErrProofNotFound
	}
	return &store.Proof{
		Address: proofs.Address,
		Denom:   proofs.Coin.Denom,
		Amount:  proofs.Coin.Amount,
		Proof:   proofs.Proof,
	}, nil
}
