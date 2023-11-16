package gorm

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"gorm.io/gorm"
)

type Account struct {
	gorm.Model
	Address       sdk.AccAddress `json:"address" gorm:"index"`
	AccountNumber int64          `json:"account_number"`
	SummaryCoins  sdk.Coins      `json:"summary_coins,omitempty"`
	Coins         sdk.Coins      `json:"coins,omitempty"`
	FrozenCoins   sdk.Coins      `json:"frozen_coins,omitempty"`
	LockedCoins   sdk.Coins      `json:"locked_coins,omitempty"`
}

type Proof struct {
	gorm.Model
	Address sdk.AccAddress `json:"address" gorm:"index"`
	Denom   string         `json:"denom" gorm:"index"`
	Proof   [][]byte       `json:"proof"`
}
