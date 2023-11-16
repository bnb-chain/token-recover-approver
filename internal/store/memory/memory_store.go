package memory

import (
	"github.com/cosmos/cosmos-sdk/types"

	"github.com/bnb-chain/airdrop-service/internal/store"
	"github.com/bnb-chain/airdrop-service/pkg/util"
)

var _ store.Store = (*MemoryStore)(nil)

func NewMemoryStore(accountsPath, proofsPath string) (*MemoryStore, error) {
	errChan := make(chan error, 1)
	defer close(errChan)
	// load accounts
	stream := util.NewJSONStream(func() any {
		return &Account{}
	})

	accounts := make(map[string]*Account)
	go func() {
		for data := range stream.Watch() {
			if data.Error != nil {
				errChan <- data.Error
				return
			}
			acc := data.Data.(*Account)
			accounts[acc.Address.String()] = acc
		}
		errChan <- nil
	}()
	stream.Start(accountsPath)
	err := <-errChan
	if err != nil {
		return nil, err
	}
	// load proofs
	stream = util.NewJSONStream(func() any {
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
	err = <-errChan
	if err != nil {
		return nil, err
	}

	return &MemoryStore{
		accounts: accounts,
		proofs:   proofs,
	}, nil
}

// MemoryStore implements store.Store.
type MemoryStore struct {
	accounts map[string]*Account
	proofs   map[string]*Proof // address:index:symbol -> proofs
}

// GetAccountByAddress implements store.Store.
func (ss *MemoryStore) GetAccountByAddress(address types.AccAddress) (*store.Account, error) {
	acc, exist := ss.accounts[address.String()]
	if !exist {
		return nil, ErrAccountNotFound
	}

	return &store.Account{
		Address:       acc.Address,
		AccountNumber: acc.AccountNumber,
		SummaryCoins:  acc.SummaryCoins,
		Coins:         acc.Coins,
		FrozenCoins:   acc.FrozenCoins,
		LockedCoins:   acc.LockedCoins,
	}, nil
}

// GetAccountProofs implements store.Store.
func (ss *MemoryStore) GetAccountAssetProof(address types.AccAddress, symbol string) ([][]byte, error) {
	index := address.String() + ":" + symbol
	proofs, exist := ss.proofs[index]
	if !exist {
		return nil, ErrProofNotFound
	}
	return proofs.Proof, nil
}
