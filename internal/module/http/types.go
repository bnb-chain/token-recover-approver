package http

import (
	"github.com/ethereum/go-ethereum/common"
)

type GetClaimApprovalRequest struct {
	TokenSymbol    string
	OwnerSignature []byte
	ClaimAddress   common.Address
}

type GetClaimApprovalResponse struct {
	Amount     int64
	Node       []byte
	PrefixNode []byte
	SuffixNode []byte
	Proofs     [][]byte
	Approval   []byte
}

type GetRegisterTokenApprovalRequest struct {
	TokenSymbol     string
	OwnerSignature  []byte
	RegisterAddress common.Address
}

type GetRegisterTokenApprovalResponse struct {
	Amount   int64
	Approval []byte
}
