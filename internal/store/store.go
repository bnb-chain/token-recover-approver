package store

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Store interface {
	GetStateRoot() (stateRoot string, err error)
	GetAssetBySymbol(symbol string) (asset *Asset, err error)
	GetAccountByAddress(address sdk.AccAddress) (account *Account, err error)
	GetAccountProofs(address sdk.AccAddress) (proofs []string, err error)
}
