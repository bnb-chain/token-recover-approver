package store

import (
	"bytes"
	"errors"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	// ErrInvalidToken is returned when the token is not found in the merkle tree
	ErrInvalidToken = errors.New("invalid token")
)

// Proof is a merkle proof of an account
type Proof struct {
	Address sdk.AccAddress `json:"address"`
	Denom   string         `json:"denom"`
	Amount  int64          `json:"amount"`
	Proof   [][]byte       `json:"proof"`
}

// Serialize implements merkle tree data Serialize method.
func (p *Proof) Serialize() ([]byte, error) {
	if p.Amount == 0 {
		return nil, ErrInvalidToken
	}

	var symbol [32]byte
	copy(symbol[:], p.Denom)
	buf := new(bytes.Buffer)
	if _, err := buf.Write(p.Address.Bytes()); err != nil {
		return nil, err
	}
	if _, err := buf.Write(symbol[:]); err != nil {
		return nil, err
	}
	if _, err := buf.Write(big.NewInt(p.Amount).FillBytes(make([]byte, 32))); err != nil {
		return nil, err
	}
	return crypto.Keccak256Hash(
		buf.Bytes(),
	).Bytes(), nil
}
