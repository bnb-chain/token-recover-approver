package store

import (
	"errors"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	// ErrInvalidTokenIndex is the error returned when the token index is invalid.
	ErrInvalidTokenIndex = errors.New("invalid token index")
)

// Account is an account.
type Account struct {
	Address       sdk.AccAddress `json:"address"`
	AccountNumber int64          `json:"account_number"`
	SummaryCoins  sdk.Coins      `json:"summary_coins,omitempty"`
	Coins         sdk.Coins      `json:"coins,omitempty"`
	FrozenCoins   sdk.Coins      `json:"frozen_coins,omitempty"`
	LockedCoins   sdk.Coins      `json:"locked_coins,omitempty"`
}

// Serialize implements merkle tree data Serialize method.
func (acc *Account) Serialize(tokenSymbol string) ([]byte, error) {
	if acc.SummaryCoins.AmountOf(tokenSymbol) == 0 {
		return nil, ErrInvalidTokenIndex
	}

	var symbol [32]byte
	copy(symbol[:], tokenSymbol)
	return crypto.Keccak256Hash(
		acc.Address.Bytes(),
		symbol[:],
		big.NewInt(acc.SummaryCoins.AmountOf(tokenSymbol)).FillBytes(make([]byte, 32)),
	).Bytes(), nil
}

// Proofs is a list of account to merkle proof
type Proofs []*Proof

// Proof is a merkle proof of an account
type Proof struct {
	Address sdk.AccAddress `json:"address"`
	Index   int64          `json:"index"`
	Coin    sdk.Coin       `json:"coin"`
	Proof   [][]byte       `json:"proof"`
}
