package approval

import (
	"encoding/hex"
	"errors"
	"math/big"

	"github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/crypto"

	// "github.com/tendermint/tendermint/crypto/secp256k1"

	"github.com/bnb-chain/airdrop-service/internal/config"
	"github.com/bnb-chain/airdrop-service/internal/store"
	"github.com/bnb-chain/airdrop-service/pkg/keymanager"
)

const (
	RequestTypeClaim         = "claim"
	RequestTypeRegisterToken = "register_token"
)

type ApprovalService struct {
	config *config.Config
	km     keymanager.KeyManager
	store  store.Store
}

func NewApprovalService(config *config.Config, km keymanager.KeyManager, store store.Store) *ApprovalService {
	return &ApprovalService{km: km, store: store}
}

func (svc *ApprovalService) GetClaimApproval(req *GetClaimApprovalRequest) (resp *GetClaimApprovalResponse, err error) {
	// Verify owner signature
	ownerSignature, err := hex.DecodeString(req.OwnerSignature)
	if err != nil {
		return nil, err
	}
	var tokenSymbolBytes [32]byte
	copy(tokenSymbolBytes[:], []byte(req.TokenSymbol))
	msg := crypto.Keccak256([]byte(RequestTypeClaim), []byte(svc.config.ChainID), tokenSymbolBytes[:], req.ClaimAddress[:])
	ownerAddr, err := svc.RecoverAddressFromTmSig(msg, ownerSignature)
	if err != nil {
		return nil, err
	}
	// Get Merkle Proofs and Node
	proofs, err := svc.store.GetAccountProofs(ownerAddr)
	if err != nil {
		return nil, err
	}
	account, err := svc.store.GetAccountByAddress(ownerAddr)
	if err != nil {
		return nil, err
	}
	nodeBytes, err := account.Serialize()
	if err != nil {
		return nil, err
	}
	prefixNode, suffixNode := account.GetPrefixSuffixNode(req.TokenSymbol)
	approvalSignature := crypto.Keccak256([]byte(svc.config.ChainID), tokenSymbolBytes[:], req.ClaimAddress[:], ownerSignature, nodeBytes)
	return &GetClaimApprovalResponse{
		Amount:            big.NewInt(account.SummaryCoins.AmountOf(req.TokenSymbol)),
		PrefixNode:        hex.EncodeToString(prefixNode),
		SuffixNode:        hex.EncodeToString(suffixNode),
		Proofs:            proofs,
		ApprovalSignature: hex.EncodeToString(approvalSignature),
	}, nil
}

func (svc *ApprovalService) GetRegisterTokenApproval(req *GetRegisterTokenApprovalRequest) (resp *GetRegisterTokenApprovalResponse, err error) {
	// Verify owner signature
	ownerSignature, err := hex.DecodeString(req.OwnerSignature)
	if err != nil {
		return nil, err
	}
	var tokenSymbolBytes [32]byte
	copy(tokenSymbolBytes[:], []byte(req.TokenSymbol))
	msg := crypto.Keccak256([]byte(RequestTypeRegisterToken), []byte(svc.config.ChainID), tokenSymbolBytes[:], req.RegisterAddress[:])
	ownerAddr, err := svc.RecoverAddressFromTmSig(msg, ownerSignature)
	if err != nil {
		return nil, err
	}
	// Get Asset Owner
	asset, err := svc.store.GetAssetBySymbol(req.TokenSymbol)
	if err != nil {
		return nil, err
	}
	if !asset.Owner.Equals(ownerAddr) {
		return nil, errors.New("asset owner not match")
	}
	amount := big.NewInt(asset.Amount)
	approvalSignature := crypto.Keccak256([]byte(svc.config.ChainID), tokenSymbolBytes[:], req.RegisterAddress[:], ownerSignature, amount.Bytes())
	return &GetRegisterTokenApprovalResponse{
		Amount:   amount,
		Approval: hex.EncodeToString(approvalSignature),
	}, nil
}

func (svc *ApprovalService) RecoverAddressFromTmSig(msg []byte, sig []byte) (types.AccAddress, error) {
	// TODO: Implement

	// pubKey, err := secp256k1.RecoverPubkey(msg, sig)
	// if err != nil {
	// 	return nil, err
	// }

	return nil, nil
}
