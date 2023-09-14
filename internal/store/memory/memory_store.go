package memory

import (
	"encoding/json"
	"io"
	"os"

	"github.com/cosmos/cosmos-sdk/types"

	"github.com/bnb-chain/airdrop-service/internal/store"
	"github.com/bnb-chain/airdrop-service/pkg/util"
)

var _ store.Store = (*MemoryStore)(nil)

func NewMemoryStore(chainID, stateRootPath, assetsPath, accountsPath, proofsPath string) (*MemoryStore, error) {
	// load state root
	var stateRoot StateRoot
	stateRootFile, err := os.Open(stateRootPath)
	if err != nil {
		return nil, err
	}
	defer stateRootFile.Close()
	buf, err := io.ReadAll(stateRootFile)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(buf, &stateRoot)
	if err != nil {
		return nil, err
	}

	// load assets
	var assets Assets
	assetsFile, err := os.Open(assetsPath)
	if err != nil {
		return nil, err
	}
	defer assetsFile.Close()
	buf, err = io.ReadAll(assetsFile)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(buf, &assets)
	if err != nil {
		return nil, err
	}

	errChan := make(chan error, 1)
	defer close(errChan)
	// load accounts
	stream := util.NewJSONStream(func() any {
		return &Account{}
	})

	accounts := map[string]*Account{}
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
	err = <-errChan
	if err != nil {
		return nil, err
	}
	// load proofs
	stream = util.NewJSONStream(func() any {
		return &Proof{}
	})

	proofs := map[string]*Proof{}
	go func() {
		for data := range stream.Watch() {
			if data.Error != nil {
				errChan <- data.Error
				return
			}
			proof := data.Data.(*Proof)
			proofs[proof.Address.String()] = proof
		}
		errChan <- nil
	}()
	stream.Start(proofsPath)
	err = <-errChan
	if err != nil {
		return nil, err
	}

	return &MemoryStore{
		stateRoot: stateRoot,
		assets:    assets,
		accounts:  accounts,
		proofs:    proofs,
	}, nil
}

// MemoryStore implements store.Store.
type MemoryStore struct {
	stateRoot StateRoot
	assets    Assets
	accounts  map[string]*Account
	proofs    map[string]*Proof
}

// GetAccountByAddress implements store.Store.
func (*MemoryStore) GetAccountByAddress(address types.AccAddress) (account *store.Account, err error) {
	panic("unimplemented")
}

// GetAccountProofs implements store.Store.
func (*MemoryStore) GetAccountProofs(address types.AccAddress) (proofs []string, err error) {
	panic("unimplemented")
}

// GetAssetBySymbol implements store.Store.
func (*MemoryStore) GetAssetBySymbol(symbol string) (asset *store.Asset, err error) {
	panic("unimplemented")
}

// GetStateRoot implements store.Store.
func (*MemoryStore) GetStateRoot() (stateRoot []byte, err error) {
	panic("unimplemented")
}
