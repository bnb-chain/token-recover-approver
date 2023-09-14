package memory

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func init() {
	sdk.GetConfig().GetBech32AccountAddrPrefix()
}

type StateRoot struct {
	StateRoot string `json:"state_root"`
}

type Account struct {
	Address       sdk.AccAddress `json:"address"`
	AccountNumber int64          `json:"account_number"`
	SummaryCoins  sdk.Coins      `json:"summary_coins,omitempty"`
	Coins         sdk.Coins      `json:"coins,omitempty"`
	FrozenCoins   sdk.Coins      `json:"frozen_coins,omitempty"`
	LockedCoins   sdk.Coins      `json:"locked_coins,omitempty"`
}

type Assets map[string]*Asset

type Asset struct {
	Owner  sdk.AccAddress `json:"owner,omitempty"`
	Amount int64          `json:"amount"`
}

type Proof struct {
	Address sdk.AccAddress `json:"address"`
	Proof   []string       `json:"proofs"`
}
