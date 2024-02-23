package store

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Store interface {
	GetAccountAssetProof(address sdk.AccAddress, symbol string) (proofs *Proof, err error)
	CountAccountAssetProofs() (count int64, err error)
	Close() error
}
