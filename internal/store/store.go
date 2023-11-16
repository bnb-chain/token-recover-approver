package store

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Store interface {
	GetAccountByAddress(address sdk.AccAddress) (account *Account, err error)
	GetAccountAssetProof(address sdk.AccAddress, symbol string) (proofs [][]byte, err error)
}
