package gorm

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"gorm.io/gorm"
)

type Account struct {
	gorm.Model
	Address       sdk.AccAddress `json:"address" gorm:"index"`
	AccountNumber int64          `json:"account_number"`
	Coins         sdk.Coins      `json:"coins,omitempty"`
}

type Proof struct {
	gorm.Model
	Address sdk.AccAddress `json:"address" gorm:"index"`
	Denom   string         `json:"denom" gorm:"index"`
	Proof   string         `json:"proof"`
}
