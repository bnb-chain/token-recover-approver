package gorm

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"gorm.io/gorm"
)

type Proof struct {
	gorm.Model
	Address sdk.AccAddress `json:"address" gorm:"index"`
	Denom   string         `json:"denom" gorm:"index"`
	Proof   string         `json:"proof" gorm:"type:text"` // hex encoded
	Amount  int64          `json:"amount" gorm:"index"`
}
