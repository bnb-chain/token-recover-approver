package goleveldb

import (
	"errors"

	"github.com/cosmos/cosmos-sdk/types"

	"github.com/bnb-chain/airdrop-service/internal/store"
)

var _ store.Store = (*KVStore)(nil)

func NewKVStore() (*KVStore, error) {
	return &KVStore{}, errors.New("unimplemented")
}

// KVStore implements store.Store.
type KVStore struct {
}

// GetAccountByAddress implements store.Store.
func (*KVStore) GetAccountByAddress(address types.AccAddress) (account *store.Account, err error) {
	panic("unimplemented")
}

// GetAccountProofs implements store.Store.
func (*KVStore) GetAccountProofs(address types.AccAddress) (proofs []string, err error) {
	panic("unimplemented")
}

// GetAssetBySymbol implements store.Store.
func (*KVStore) GetAssetBySymbol(symbol string) (asset *store.Asset, err error) {
	panic("unimplemented")
}

// GetStateRoot implements store.Store.
func (*KVStore) GetStateRoot() (stateRoot string, err error) {
	panic("unimplemented")
}
