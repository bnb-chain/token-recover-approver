package approval

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/ecdsa"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/crypto/ripemd160"

	"github.com/bnb-chain/node/plugins/airdrop"

	"github.com/bnb-chain/airdrop-service/internal/config"
	"github.com/bnb-chain/airdrop-service/internal/store"
	"github.com/bnb-chain/airdrop-service/pkg/keymanager"
)

const (
	RequestTypeClaim         = "claim"
	RequestTypeRegisterToken = "register_token"
)

type ApprovalService struct {
	config           *config.Config
	km               keymanager.KeyManager
	store            store.Store
	accountWhileList map[string]struct{}
}

func NewApprovalService(config *config.Config, km keymanager.KeyManager, store store.Store) *ApprovalService {
	accountWhileList := make(map[string]struct{})
	for _, addr := range config.AccountWhileList {
		accountWhileList[addr] = struct{}{}
	}
	return &ApprovalService{km: km, store: store, config: config, accountWhileList: accountWhileList}
}

func (svc *ApprovalService) checkWhileList(acc types.AccAddress) bool {
	if len(svc.accountWhileList) == 0 {
		return true
	}

	_, ok := svc.accountWhileList[acc.String()]
	return ok
}

func (svc *ApprovalService) GetClaimApproval(req *GetClaimApprovalRequest) (resp *GetClaimApprovalResponse, err error) {
	ownerPubKeyBytes, err := hex.DecodeString(req.OwnerPubKey)
	if err != nil {
		return nil, err
	}
	ownerAddr, err := svc.getAddressFromPubKey(ownerPubKeyBytes)
	if err != nil {
		return nil, err
	}
	ownerSignature, err := hex.DecodeString(req.OwnerSignature)
	if err != nil {
		return nil, err
	}

	// Check While List
	if !svc.checkWhileList(ownerAddr) {
		return nil, errors.New("address is not in while list")
	}

	// Get Merkle Proofs and Node
	proofs, err := svc.store.GetAccountAssetProofs(ownerAddr, req.TokenSymbol, int64(req.TokenIndex))
	if err != nil {
		return nil, err
	}
	account, err := svc.store.GetAccountByAddress(ownerAddr)
	if err != nil {
		return nil, err
	}

	// Check if token amount is zero
	if account.SummaryCoins[req.TokenIndex].Amount == 0 {
		return nil, errors.New("token amount is zero")
	}

	// Verify user signature
	approvalMsg := airdrop.NewAirdropApprovalMsg(req.TokenIndex, req.TokenSymbol, uint64(account.SummaryCoins[req.TokenIndex].Amount), req.ClaimAddress.String())
	hash := sha256.New()
	hash.Write(approvalMsg.GetSignBytes())
	msgHash := hash.Sum(nil)
	err = svc.verifyTmSignature(ownerPubKeyBytes, ownerSignature, msgHash)
	if err != nil {
		return nil, err
	}

	nodeBytes, err := account.Serialize(uint(req.TokenIndex))
	if err != nil {
		return nil, err
	}

	var tokenSymbolBytes [32]byte
	copy(tokenSymbolBytes[:], []byte(req.TokenSymbol))

	approvalSignature := crypto.Keccak256([]byte(svc.config.ChainID), tokenSymbolBytes[:], req.ClaimAddress[:], ownerSignature, nodeBytes)
	return &GetClaimApprovalResponse{
		Amount:            big.NewInt(account.SummaryCoins.AmountOf(req.TokenSymbol)),
		Proofs:            proofs,
		ApprovalSignature: hex.EncodeToString(approvalSignature),
	}, nil
}

func (svc *ApprovalService) verifyTmSignature(pubkey, signatureStr, msgHash []byte) error {
	pubKey, err := btcec.ParsePubKey(pubkey)
	if err != nil {
		return err
	}

	r, s, err := svc.signatureFromBytes(signatureStr)
	if err != nil {
		return err
	}
	signature := ecdsa.NewSignature(r, s)

	// Reject malleable signatures. libsecp256k1 does this check but btcec doesn't.
	if s.IsOverHalfOrder() {
		return fmt.Errorf("invalid signature")
	}

	// Verify the signature.
	if !signature.Verify(msgHash, pubKey) {
		return fmt.Errorf("invalid signature")
	}

	return nil
}

func (svc ApprovalService) getAddressFromPubKey(pubkey []byte) (types.AccAddress, error) {
	pubKey, err := btcec.ParsePubKey(pubkey)
	if err != nil {
		return nil, err
	}
	hasherSHA256 := sha256.New()
	_, _ = hasherSHA256.Write(pubKey.SerializeCompressed()) // does not error
	sha := hasherSHA256.Sum(nil)

	hasherRIPEMD160 := ripemd160.New()
	_, _ = hasherRIPEMD160.Write(sha) // does not error
	return hasherRIPEMD160.Sum(nil), nil
}

// Read Signature struct from R || S. Caller needs to ensure
// that len(sigStr) == 64.
func (svc ApprovalService) signatureFromBytes(sigStr []byte) (*btcec.ModNScalar, *btcec.ModNScalar, error) {
	var r, s btcec.ModNScalar
	if r.SetByteSlice(sigStr[:32]) {
		return nil, nil, fmt.Errorf("invalid R field")
	}
	if s.SetByteSlice(sigStr[32:]) {
		return nil, nil, fmt.Errorf("invalid S field")
	}
	return &r, &s, nil
}
