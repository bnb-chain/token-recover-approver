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

// Account is an  account.
type Account struct {
	Address       sdk.AccAddress `json:"address"`
	AccountNumber int64          `json:"account_number"`
	SummaryCoins  sdk.Coins      `json:"summary_coins,omitempty"`
	Coins         sdk.Coins      `json:"coins,omitempty"`
	FrozenCoins   sdk.Coins      `json:"frozen_coins,omitempty"`
	LockedCoins   sdk.Coins      `json:"locked_coins,omitempty"`
}

// Serialize implements merkle tree data Serialize method.
func (acc *Account) Serialize(tokenIndex uint) ([]byte, error) {
	if tokenIndex >= uint(acc.SummaryCoins.Len()) {
		return nil, ErrInvalidTokenIndex
	}

	if acc.SummaryCoins[tokenIndex].Amount == 0 {
		return nil, ErrInvalidTokenIndex
	}

	var symbol [32]byte
	copy(symbol[:], acc.SummaryCoins[tokenIndex].Denom)
	return crypto.Keccak256Hash(
		acc.Address.Bytes(),
		big.NewInt(int64(tokenIndex)).FillBytes(make([]byte, 32)),
		symbol[:],
		big.NewInt(acc.SummaryCoins[tokenIndex].Amount).FillBytes(make([]byte, 32)),
	).Bytes(), nil
}

// Assets is a map of asset name to amount
type Assets map[string]*Asset

type Asset struct {
	Owner  sdk.AccAddress `json:"owner,omitempty"`
	Amount int64          `json:"amount"`
}

// Proofs is a map of account address to merkle proof
type Proofs map[string][]string

// WorldState is the world state of the store.
type WorldState struct {
	ChainID     string       `json:"chain_id"`
	BlockHeight int64        `json:"block_height"`
	CommitID    sdk.CommitID `json:"commit_id"`
	Accounts    []*Account   `json:"accounts"`
	Assets      Assets       `json:"assets"`
	StateRoot   string       `json:"state_root"`
	Proofs      Proofs       `json:"proofs"`
}
