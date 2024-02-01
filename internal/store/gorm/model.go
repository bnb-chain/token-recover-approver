package gorm

import (
	"gorm.io/gorm"
)

type Proof struct {
	gorm.Model
	Address string `json:"address" gorm:"index,type:varchar(42)"` // hex encoded
	Denom   string `json:"denom" gorm:"index"`
	Proof   string `json:"proof" gorm:"type:text"` // hex encoded
	Amount  int64  `json:"amount"`
}
