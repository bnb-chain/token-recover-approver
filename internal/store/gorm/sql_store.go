package gorm

import (
	"errors"

	"github.com/cosmos/cosmos-sdk/types"

	"github.com/bnb-chain/airdrop-service/internal/store"
)

var _ store.Store = (*SQLStore)(nil)

func NewSQLStore() (*SQLStore, error) {
	return &SQLStore{}, errors.New("unimplemented")
}

// SQLStore implements store.Store.
type SQLStore struct {
}

// GetAccountByAddress implements store.Store.
func (*SQLStore) GetAccountByAddress(address types.AccAddress) (account *store.Account, err error) {
	panic("unimplemented")
}

// GetAccountProofs implements store.Store.
func (*SQLStore) GetAccountProofs(address types.AccAddress) (proofs []string, err error) {
	panic("unimplemented")
}

// GetAssetBySymbol implements store.Store.
func (*SQLStore) GetAssetBySymbol(symbol string) (asset *store.Asset, err error) {
	panic("unimplemented")
}

// GetStateRoot implements store.Store.
func (*SQLStore) GetStateRoot() (stateRoot string, err error) {
	panic("unimplemented")
}
