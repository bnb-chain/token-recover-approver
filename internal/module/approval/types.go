package approval

import (
	"encoding/json"
	"math/big"

	"github.com/pkg/errors"

	"github.com/bnb-chain/token-recover-approver/pkg/util"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

const (
	PubKeyLength    = 33
	SignatureLength = 64
)

type GetTokenRecoverApprovalRequest struct {
	TokenSymbol    string         `json:"token_symbol" validate:"required"`
	OwnerPubKey    string         `json:"owner_pub_key" validate:"required"`
	OwnerSignature string         `json:"owner_signature" validate:"required"`
	ClaimAddress   common.Address `json:"claim_address" validate:"required"`
}

func (req *GetTokenRecoverApprovalRequest) Validate() error {
	if (req.ClaimAddress == common.Address{}) {
		return errors.New("claim address is empty")
	}

	pubKey, err := hexutil.Decode(req.OwnerPubKey)
	if err != nil {
		return errors.Wrap(err, "decode owner public key")
	}

	if len(pubKey) != PubKeyLength {
		return errors.New("invalid owner public key")
	}

	signature, err := hexutil.Decode(req.OwnerSignature)
	if err != nil {
		return errors.Wrap(err, "decode owner signature")
	}

	if len(signature) != SignatureLength {
		return errors.New("invalid owner signature")
	}

	return nil
}

type GetTokenRecoverApprovalResponse struct {
	Amount            *big.Int `json:"amount"`
	Proofs            [][]byte `json:"proofs"`
	ApprovalSignature []byte   `json:"approval_signature"`
}

func (resp *GetTokenRecoverApprovalResponse) MarshalJSON() ([]byte, error) {
	type aliasGetTokenRecoverApprovalResponse struct {
		Amount            *big.Int `json:"amount"`
		Proofs            []string `json:"proofs"`
		ApprovalSignature string   `json:"approval_signature"`
	}
	return json.Marshal(&aliasGetTokenRecoverApprovalResponse{
		Amount:            resp.Amount,
		Proofs:            util.EncodeBytesArrayToHex(resp.Proofs),
		ApprovalSignature: hexutil.Encode(resp.ApprovalSignature),
	})
}
