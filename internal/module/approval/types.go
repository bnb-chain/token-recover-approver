package approval

import (
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type GetClaimApprovalRequest struct {
	TokenIndex     uint64         `json:"token_index" validate:"required"`
	TokenSymbol    string         `json:"token_symbol" validate:"required"`
	OwnerPubKey    string         `json:"owner_pub_key" validate:"required"`
	OwnerSignature string         `json:"owner_signature" validate:"required"`
	ClaimAddress   common.Address `json:"claim_address" validate:"required"`
}

func (req *GetClaimApprovalRequest) Validate() error {
	if (req.ClaimAddress == common.Address{}) {
		return errors.New("claim address is empty")
	}

	return nil
}

type GetClaimApprovalResponse struct {
	Amount            *big.Int `json:"amount"`
	Proofs            []string `json:"proofs"`
	ApprovalSignature string   `json:"approval_signature"`
}
