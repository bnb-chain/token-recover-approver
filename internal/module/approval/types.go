package approval

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type GetClaimApprovalRequest struct {
	TokenSymbol    string         `json:"token_symbol" validate:"required"`
	OwnerSignature string         `json:"owner_signature" validate:"required"`
	ClaimAddress   common.Address `json:"claim_address" validate:"required"`
}

type GetClaimApprovalResponse struct {
	Amount            *big.Int `json:"amount"`
	PrefixNode        string   `json:"prefix_node"`
	SuffixNode        string   `json:"suffix_node"`
	Proofs            []string `json:"proofs"`
	ApprovalSignature string   `json:"approval_signature"`
}

type GetRegisterTokenApprovalRequest struct {
	TokenSymbol     string         `json:"token_symbol" validate:"required"`
	OwnerSignature  string         `json:"owner_signature" validate:"required"`
	RegisterAddress common.Address `json:"register_address" validate:"required"`
}

type GetRegisterTokenApprovalResponse struct {
	Amount   *big.Int `json:"amount"`
	Approval string   `json:"approval"`
}
