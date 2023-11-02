package gorm

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"gorm.io/gorm"
)

func init() {
	sdk.GetConfig().GetBech32AccountAddrPrefix()
}

type StateRoot struct {
	gorm.Model
	StateRoot string `json:"state_root"`
}

type Account struct {
	gorm.Model
	Address       sdk.AccAddress `json:"address"`
	AccountNumber int64          `json:"account_number"`
	SummaryCoins  sdk.Coins      `json:"summary_coins,omitempty"`
	Coins         sdk.Coins      `json:"coins,omitempty"`
	FrozenCoins   sdk.Coins      `json:"frozen_coins,omitempty"`
	LockedCoins   sdk.Coins      `json:"locked_coins,omitempty"`
}

type Asset struct {
	gorm.Model
	Owner  sdk.AccAddress `json:"owner,omitempty"`
	Amount int64          `json:"amount"`
}

type Proof struct {
	gorm.Model
	Address sdk.AccAddress `json:"address"`
	Index   int64          `json:"index"`
	Coin    sdk.Coin       `json:"coin"`
	Proof   []string       `json:"proof"`
}
