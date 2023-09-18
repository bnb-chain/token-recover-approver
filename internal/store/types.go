package store

import (
	"bytes"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/crypto"
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
func (acc *Account) Serialize() ([]byte, error) {
	coinBytes := bytes.NewBuffer(nil)
	for index, coin := range acc.SummaryCoins {
		coinBytes.Write(tokenHash(int64(index), coin.Denom, coin.Amount))
	}
	return crypto.Keccak256Hash(
		acc.Address.Bytes(),
		big.NewInt(acc.AccountNumber).Bytes(),
		coinBytes.Bytes(),
	).Bytes(), nil
}

func (acc *Account) GetPrefixSuffixNode(symbol string) ([]byte, []byte) {
	prefixBytes := bytes.NewBuffer(nil)
	suffixBytes := bytes.NewBuffer(nil)

	prefixBytes.Write(acc.Address.Bytes())
	prefixBytes.Write(big.NewInt(acc.AccountNumber).Bytes())

	isSplit := false
	for index, coin := range acc.SummaryCoins {
		if coin.Denom == symbol {
			isSplit = true
			continue
		}

		if !isSplit {
			prefixBytes.Write(tokenHash(int64(index), coin.Denom, coin.Amount))
		} else {
			suffixBytes.Write(tokenHash(int64(index), coin.Denom, coin.Amount))
		}
	}

	return prefixBytes.Bytes(), suffixBytes.Bytes()
}

func tokenHash(index int64, denom string, amount int64) []byte {
	var symbol [32]byte
	copy(symbol[:], denom)

	return crypto.Keccak256Hash(
		big.NewInt(index).Bytes(),
		symbol[:],
		big.NewInt(amount).Bytes()).Bytes()
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
