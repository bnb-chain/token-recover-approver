package memory

import (
	"encoding/json"
	"io"
	"os"
	"strconv"

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
	err = <-errChan
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
			tokenIndexStr := strconv.Itoa(int(proof.Index))
			index := proof.Address.String() + ":" + tokenIndexStr + ":" + proof.Coin.Denom
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
	proofs    map[string]*Proof // address:index:symbol -> proofs
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
func (ss *MemoryStore) GetAccountAssetProofs(address types.AccAddress, symbol string, tokenIndex int64) ([]string, error) {
	tokenIndexStr := strconv.Itoa(int(tokenIndex))
	index := address.String() + ":" + tokenIndexStr + ":" + symbol
	proofs, exist := ss.proofs[index]
	if !exist {
		return nil, ErrProofNotFound
	}
	return proofs.Proof, nil
}

// GetAssetBySymbol implements store.Store.
func (ss *MemoryStore) GetAssetBySymbol(symbol string) (*store.Asset, error) {
	asset, exist := ss.assets[symbol]
	if !exist {
		return nil, ErrAssetNotFound
	}

	return &store.Asset{
		Owner:  asset.Owner,
		Amount: asset.Amount,
	}, nil
}

// GetStateRoot implements store.Store.
func (ss *MemoryStore) GetStateRoot() (stateRoot string, err error) {
	return ss.stateRoot.StateRoot, nil
}
